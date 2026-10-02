package backup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go-server/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type DriveFile struct {
	ID           string            `json:"id"`
	Size         int64             `json:"size,string"`
	MD5          string            `json:"md5Checksum"`
	Created      string            `json:"createdTime"`
	Parents      []string          `json:"parents"`
	Properties   map[string]string `json:"appProperties"`
	MimeType     string            `json:"mimeType"`
	Trashed      bool              `json:"trashed"`
	Capabilities struct {
		CanAddChildren bool `json:"canAddChildren"`
	} `json:"capabilities"`
}
type Cloud interface {
	Preflight(context.Context) error
	Upload(context.Context, *Run, string) (string, error)
	Verify(context.Context, *Run) (string, error)
	Retain(context.Context, *Run) (int, error)
}
type Drive struct {
	client               *http.Client
	base, upload, folder string
}

type driveHTTPError struct{ status int }

func (e *driveHTTPError) Error() string { return fmt.Sprintf("Drive HTTP %d", e.status) }

func NewDrive(cfg config.BackupConfig) *Drive {
	folder, _ := config.DriveFolderID(cfg.Folder)
	oauth := &oauth2.Config{ClientID: cfg.ClientID, ClientSecret: cfg.ClientSecret, Endpoint: google.Endpoint}
	// Refresh requests must also be bounded independently of a backup's total timeout.
	tokenCtx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Timeout: 30 * time.Second})
	client := oauth2.NewClient(tokenCtx, oauth.TokenSource(tokenCtx, &oauth2.Token{RefreshToken: cfg.RefreshToken}))
	client.Timeout = 2 * time.Minute
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	return &Drive{client: client, base: "https://www.googleapis.com/drive/v3", upload: "https://www.googleapis.com/upload/drive/v3", folder: folder}
}

func pause(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
func retryDelay(h http.Header, attempt int) time.Duration {
	d := time.Second * time.Duration(1<<attempt)
	if seconds, err := strconv.Atoi(h.Get("Retry-After")); err == nil && seconds > 0 {
		d = time.Duration(seconds) * time.Second
	}
	if d > 30*time.Second {
		d = 30 * time.Second
	}
	return d
}

// request retries only replayable operations. Bodies are bounded JSON, never archives.
func (d *Drive) request(ctx context.Context, method, endpoint string, body any, out any) error {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	for attempt := 0; attempt < 5; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		res, err := d.client.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if attempt == 4 {
				return errors.New("Drive request unavailable")
			}
			if err = pause(ctx, retryDelay(nil, attempt)); err != nil {
				return err
			}
			continue
		}
		if res.StatusCode == 429 || res.StatusCode >= 500 {
			_ = res.Body.Close()
			if attempt == 4 {
				return errors.New("Drive retry limit exceeded")
			}
			if err = pause(ctx, retryDelay(res.Header, attempt)); err != nil {
				return err
			}
			continue
		}
		defer func() { _ = res.Body.Close() }()
		if method == http.MethodDelete && res.StatusCode == 404 {
			return nil
		}
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return &driveHTTPError{status: res.StatusCode}
		}
		if out != nil {
			return json.NewDecoder(io.LimitReader(res.Body, 8<<20)).Decode(out)
		}
		return nil
	}
	return errors.New("Drive retry limit exceeded")
}
func (d *Drive) file(ctx context.Context, id string) (*DriveFile, error) {
	f := new(DriveFile)
	err := d.request(ctx, http.MethodGet, d.base+"/files/"+url.PathEscape(id)+"?fields=id,size,md5Checksum,createdTime,parents,appProperties,mimeType,trashed,capabilities", nil, f)
	return f, err
}
func (d *Drive) Preflight(ctx context.Context) error {
	f, err := d.file(ctx, d.folder)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var response *driveHTTPError
		if errors.As(err, &response) {
			switch response.status {
			case http.StatusUnauthorized:
				return preflightFailure("drive_auth", "Drive rejected authentication; check the backup OAuth credentials and refresh token")
			case http.StatusForbidden:
				return preflightFailure("drive_forbidden", "Drive denied folder access; check OAuth scopes, folder permissions, and Drive API enablement")
			case http.StatusNotFound:
				return preflightFailure("drive_folder_missing", "Drive folder was not found or is inaccessible to the backup account; check BACKUP_GDRIVE_FOLDER")
			}
		}
		return preflightFailure("drive_request", "Cannot inspect the Drive backup folder; check OAuth refresh credentials, network access, and Drive availability")
	}
	if f.Trashed || f.MimeType != "application/vnd.google-apps.folder" || !f.Capabilities.CanAddChildren {
		return preflightFailure("drive_folder_not_writable", "Drive destination is not a writable folder; check BACKUP_GDRIVE_FOLDER and backup account permissions")
	}
	return nil
}
func properties(r *Run) map[string]string {
	return map[string]string{"system": "apna-gate-postgres-backup-v1", "deployment": r.Deployment, "database": r.Database, "type": r.Type, "run_id": r.ID, "verified": "false", "sha256": r.SHA256}
}
func (d *Drive) Upload(ctx context.Context, r *Run, path string) (string, error) {
	var generated struct {
		IDs []string `json:"ids"`
	}
	if err := d.request(ctx, http.MethodGet, d.base+"/files/generateIds?count=1&space=drive&type=files", nil, &generated); err != nil {
		return "", err
	}
	if len(generated.IDs) != 1 {
		return "", errors.New("Drive did not allocate a file ID")
	}
	id := generated.IDs[0]
	meta := map[string]any{"id": id, "name": "postgres-" + r.Type + "-" + r.Created.UTC().Format("20060102T150405Z") + "-" + r.ID + ".dump", "parents": []string{d.folder}, "appProperties": properties(r)}
	payload, _ := json.Marshal(meta)
	var session string
	for attempt := 0; attempt < 5; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.upload+"/files?uploadType=resumable", bytes.NewReader(payload))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Upload-Content-Type", "application/octet-stream")
		req.Header.Set("X-Upload-Content-Length", strconv.FormatInt(r.ArchiveSize, 10))
		res, err := d.client.Do(req)
		if err == nil {
			_ = res.Body.Close()
			if res.StatusCode >= 200 && res.StatusCode < 300 {
				session = res.Header.Get("Location")
				break
			}
			if res.StatusCode == 409 {
				return id, nil
			} // same allocated ID; verification follows
			if res.StatusCode != 429 && res.StatusCode < 500 {
				return "", fmt.Errorf("Drive upload initialization HTTP %d", res.StatusCode)
			}
		}
		if attempt == 4 {
			return "", errors.New("Drive upload initialization failed")
		}
		if err := pause(ctx, retryDelay(nil, attempt)); err != nil {
			return "", err
		}
	}
	u, err := url.Parse(session)
	base, _ := url.Parse(d.upload)
	if err != nil || u.Host != base.Host || u.Scheme != base.Scheme || u.User != nil {
		return "", errors.New("invalid Drive upload session destination")
	}
	f, err := os.Open(path) //nolint:gosec // Path is generated under the private run directory, never from a request.
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	const chunkSize int64 = 8 << 20 // multiple of Drive's 256 KiB requirement
	var offset int64
	failures := 0
	query := false
	for offset < r.ArchiveSize || query {
		n := min(chunkSize, r.ArchiveSize-offset)
		var body io.Reader = io.NewSectionReader(f, offset, n)
		contentRange := fmt.Sprintf("bytes %d-%d/%d", offset, offset+n-1, r.ArchiveSize)
		if query {
			body = nil
			n = 0
			contentRange = fmt.Sprintf("bytes */%d", r.ArchiveSize)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, session, body)
		if err != nil {
			return "", err
		}
		req.ContentLength = n
		req.Header.Set("Content-Range", contentRange)
		req.Header.Set("Content-Type", "application/octet-stream")
		res, err := d.client.Do(req)
		if err != nil {
			failures++
			if failures >= 5 {
				return "", errors.New("Drive upload retry limit exceeded")
			}
			query = true
			if err := pause(ctx, retryDelay(nil, failures-1)); err != nil {
				return "", err
			}
			continue
		}
		_ = res.Body.Close()
		if res.StatusCode == 200 || res.StatusCode == 201 {
			return id, nil
		}
		if res.StatusCode == 308 {
			next, err := uploadedOffset(res.Header.Get("Range"))
			if err != nil || next < offset || next > r.ArchiveSize {
				return "", errors.New("invalid Drive upload offset")
			}
			if next == offset {
				failures++
			} else {
				failures = 0
			}
			if failures >= 5 {
				return "", errors.New("Drive upload made no progress")
			}
			offset = next
			query = offset == r.ArchiveSize
			continue
		}
		if res.StatusCode != 429 && res.StatusCode < 500 {
			return "", fmt.Errorf("Drive upload HTTP %d", res.StatusCode)
		}
		failures++
		if failures >= 5 {
			return "", errors.New("Drive upload retry limit exceeded")
		}
		query = true
		if err := pause(ctx, retryDelay(res.Header, failures-1)); err != nil {
			return "", err
		}
	}
	return "", errors.New("Drive upload did not complete")
}
func uploadedOffset(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	if !strings.HasPrefix(value, "bytes=0-") {
		return 0, errors.New("invalid upload range")
	}
	n, err := strconv.ParseInt(strings.TrimPrefix(value, "bytes=0-"), 10, 64)
	if err != nil || n < 0 || n == int64(^uint64(0)>>1) {
		return 0, errors.New("invalid upload range")
	}
	return n + 1, nil
}
func hasParent(f *DriveFile, folder string) bool {
	for _, p := range f.Parents {
		if p == folder {
			return true
		}
	}
	return false
}
func matches(f *DriveFile, r *Run, folder string) bool {
	if f.Trashed || !hasParent(f, folder) {
		return false
	}
	for _, key := range []string{"system", "deployment", "database", "type"} {
		if f.Properties[key] != properties(r)[key] {
			return false
		}
	}
	_, err := uuid.Parse(f.Properties["run_id"])
	return err == nil
}
func (d *Drive) Verify(ctx context.Context, r *Run) (string, error) {
	f, err := d.file(ctx, r.DriveFileID)
	if err != nil {
		return "", err
	}
	if !matches(f, r, d.folder) || f.Properties["run_id"] != r.ID || f.Properties["sha256"] != r.SHA256 || f.Size != r.ArchiveSize || !strings.EqualFold(f.MD5, r.LocalMD5) {
		return "", errors.New("uploaded archive metadata or checksum mismatch")
	}
	p := properties(r)
	p["verified"] = "true"
	p["verified_at"] = time.Now().UTC().Format(time.RFC3339)
	if err := d.request(ctx, http.MethodPatch, d.base+"/files/"+url.PathEscape(r.DriveFileID), map[string]any{"appProperties": p}, nil); err != nil {
		return "", err
	}
	return f.MD5, nil
}
func quoteDrive(s string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(s) + "'"
}
func (d *Drive) Retain(ctx context.Context, r *Run) (int, error) {
	q := quoteDrive(d.folder) + " in parents and trashed = false"
	for _, k := range []string{"system", "deployment", "database", "type"} {
		q += " and appProperties has { key = " + quoteDrive(k) + " and value = " + quoteDrive(properties(r)[k]) + " }"
	}
	q += " and appProperties has { key = 'verified' and value = 'true' }"
	var files []DriveFile
	token := ""
	for {
		params := url.Values{"q": {q}, "fields": {"nextPageToken,files(id,createdTime,parents,appProperties,trashed)"}, "pageSize": {"1000"}, "orderBy": {"createdTime desc"}}
		if token != "" {
			params.Set("pageToken", token)
		}
		var page struct {
			Files []DriveFile `json:"files"`
			Next  string      `json:"nextPageToken"`
		}
		if err := d.request(ctx, http.MethodGet, d.base+"/files?"+params.Encode(), nil, &page); err != nil {
			return 0, err
		}
		files = append(files, page.Files...)
		token = page.Next
		if token == "" {
			break
		}
	}
	days := 7
	if r.Type == "weekly" {
		days = 84
	}
	cutoff := time.Now().UTC().Add(-time.Duration(days) * 24 * time.Hour)
	newest := ""
	var newestTime time.Time
	for i := range files {
		f := &files[i]
		created, err := time.Parse(time.RFC3339, f.Created)
		if err == nil && matches(f, r, d.folder) && f.Properties["verified"] == "true" && created.After(newestTime) {
			newest = f.ID
			newestTime = created
		}
	}
	deleted := 0
	for i := range files {
		f := &files[i]
		created, err := time.Parse(time.RFC3339, f.Created)
		if err != nil || f.ID == newest || f.ID == r.DriveFileID || !created.Before(cutoff) || !matches(f, r, d.folder) || f.Properties["verified"] != "true" || len(f.Properties["sha256"]) != 64 {
			continue
		}
		if _, err := time.Parse(time.RFC3339, f.Properties["verified_at"]); err != nil {
			continue
		}
		// Re-read immediately before deletion: never act on a stale parent/tag listing.
		current, err := d.file(ctx, f.ID)
		if err != nil {
			return deleted, err
		}
		if !matches(current, r, d.folder) || current.Properties["verified"] != "true" || current.Created != f.Created || current.Properties["run_id"] != f.Properties["run_id"] || current.Properties["sha256"] != f.Properties["sha256"] || current.Properties["verified_at"] != f.Properties["verified_at"] {
			continue
		}
		if err := d.request(ctx, http.MethodDelete, d.base+"/files/"+url.PathEscape(f.ID), nil, nil); err != nil {
			return deleted, err
		}
		deleted++
	}
	return deleted, nil
}
