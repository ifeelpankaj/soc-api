package backup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func testRun() *Run {
	return &Run{ID: uuid.NewString(), Database: "main_db", Deployment: "production", Type: "daily", Created: time.Now().UTC(), Details: Details{SHA256: strings.Repeat("a", 64), LocalMD5: strings.Repeat("b", 32), ArchiveSize: 100, DriveFileID: "current"}}
}
func serveJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func TestDrivePreflightDiagnosticsRedactProviderBody(t *testing.T) {
	for _, tc := range []struct {
		status int
		code   string
	}{
		{401, "preflight_drive_auth"},
		{403, "preflight_drive_forbidden"},
		{404, "preflight_drive_folder_missing"},
		{400, "preflight_drive_request"},
		{200, "preflight_drive_folder_not_writable"},
	} {
		t.Run(tc.code, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, `{"error":"secret-provider-body","mimeType":"application/vnd.google-apps.folder","capabilities":{"canAddChildren":false}}`)
			}))
			defer server.Close()
			d := &Drive{client: server.Client(), base: server.URL, folder: "folder"}
			err := d.Preflight(context.Background())
			var diagnostic *diagnosticError
			if !errors.As(err, &diagnostic) || diagnostic.code != tc.code {
				t.Fatalf("diagnostic=%v", err)
			}
			if strings.Contains(err.Error(), "secret") {
				t.Fatal("provider body leaked")
			}
		})
	}
}

func TestResumableUploadContinuesAtConfirmedOffset(t *testing.T) {
	r := testRun()
	path := filepath.Join(t.TempDir(), "dump")
	data := bytes.Repeat([]byte("x"), 100)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	var server *httptest.Server
	puts := 0
	initializations := 0
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/files/generateIds":
			serveJSON(w, map[string]any{"ids": []string{"stable-id"}})
		case "/files":
			var meta map[string]any
			if err := json.NewDecoder(req.Body).Decode(&meta); err != nil {
				t.Error(err)
			}
			if meta["id"] != "stable-id" {
				t.Error("upload retry changed file identity")
			}
			initializations++
			if initializations == 1 {
				w.WriteHeader(503)
				return
			}
			w.Header().Set("Location", server.URL+"/session")
			w.WriteHeader(200)
		case "/session":
			puts++
			body, _ := io.ReadAll(req.Body)
			switch puts {
			case 1:
				if req.Header.Get("Content-Range") != "bytes 0-99/100" || len(body) != 100 {
					t.Error("wrong initial chunk")
				}
				w.WriteHeader(503)
			case 2:
				if req.Header.Get("Content-Range") != "bytes */100" || len(body) != 0 {
					t.Error("missing resumable status query")
				}
				w.Header().Set("Range", "bytes=0-49")
				w.WriteHeader(308)
			case 3:
				if req.Header.Get("Content-Range") != "bytes 50-99/100" || len(body) != 50 {
					t.Error("did not resume at confirmed offset")
				}
				w.WriteHeader(201)
			default:
				t.Error("unexpected chunk")
			}
		default:
			t.Errorf("unexpected path %s", req.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	d := &Drive{client: server.Client(), base: server.URL, upload: server.URL, folder: "folder"}
	id, err := d.Upload(context.Background(), r, path)
	if err != nil || id != "stable-id" || puts != 3 || initializations != 2 {
		t.Fatalf("id %q err %v puts %d", id, err, puts)
	}
}
func TestUploadRejectsExternalSessionURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/files/generateIds" {
			serveJSON(w, map[string]any{"ids": []string{"id"}})
			return
		}
		w.Header().Set("Location", "https://attacker.invalid/session")
		w.WriteHeader(200)
	}))
	defer srv.Close()
	d := &Drive{client: srv.Client(), base: srv.URL, upload: srv.URL, folder: "folder"}
	if _, err := d.Upload(context.Background(), testRun(), "unused"); err == nil {
		t.Fatal("accepted external session")
	}
}
func TestVerifyRequiresMatchingHashSizeAndIdentity(t *testing.T) {
	for _, change := range []string{"none", "size", "md5", "parent", "database", "run_id", "sha256"} {
		t.Run(change, func(t *testing.T) {
			r := testRun()
			f := DriveFile{ID: r.DriveFileID, Size: r.ArchiveSize, MD5: r.LocalMD5, Parents: []string{"folder"}, Properties: properties(r)}
			switch change {
			case "size":
				f.Size++
			case "md5":
				f.MD5 = "wrong"
			case "parent":
				f.Parents = []string{"elsewhere"}
			case "database":
				f.Properties["database"] = "another"
			case "run_id":
				f.Properties["run_id"] = uuid.NewString()
			case "sha256":
				f.Properties["sha256"] = "wrong"
			}
			marked := false
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method == "PATCH" {
					marked = true
					serveJSON(w, map[string]any{})
					return
				}
				serveJSON(w, f)
			}))
			defer srv.Close()
			d := &Drive{client: srv.Client(), base: srv.URL, folder: "folder"}
			_, err := d.Verify(context.Background(), r)
			if change == "none" {
				if err != nil || !marked {
					t.Fatal(err)
				}
			} else if err == nil || marked {
				t.Fatalf("verified bad %s", change)
			}
		})
	}
}
func TestRetentionMetadataBoundariesAndPagination(t *testing.T) {
	for _, kind := range []string{"daily", "weekly"} {
		t.Run(kind, func(t *testing.T) {
			r := testRun()
			r.Type = kind
			days := 7
			if kind == "weekly" {
				days = 84
			}
			makeFile := func(id string, age int) DriveFile {
				p := properties(r)
				p["verified"] = "true"
				p["verified_at"] = time.Now().Add(-time.Duration(age) * 24 * time.Hour).UTC().Format(time.RFC3339)
				p["run_id"] = uuid.NewString()
				return DriveFile{ID: id, Created: time.Now().Add(-time.Duration(age) * 24 * time.Hour).UTC().Format(time.RFC3339), Properties: p, Parents: []string{"folder"}}
			}
			files := []DriveFile{makeFile("current", 0), makeFile("recent", days-1), makeFile("expired", days+1), makeFile("other-folder", days+2), makeFile("other-db", days+2), makeFile("other-type", days+2), makeFile("other-system", days+2), makeFile("unverified", days+2), makeFile("missing-run", days+2), makeFile("moved-after-list", days+2)}
			files[3].Parents = []string{"other"}
			files[4].Properties["database"] = "other"
			files[5].Properties["type"] = "other"
			files[6].Properties["system"] = "other"
			files[7].Properties["verified"] = "false"
			delete(files[8].Properties, "run_id")
			deleted := []string{}
			pages := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/files" {
					q := req.URL.Query().Get("q")
					if !strings.Contains(q, "appProperties") || !strings.Contains(q, "'folder' in parents") {
						t.Error("retention query is not scoped")
					}
					pages++
					if req.URL.Query().Get("pageToken") == "" {
						serveJSON(w, map[string]any{"files": files[:5], "nextPageToken": "next"})
					} else {
						serveJSON(w, map[string]any{"files": files[5:]})
					}
					return
				}
				id := strings.TrimPrefix(req.URL.Path, "/files/")
				if req.Method == "DELETE" {
					deleted = append(deleted, id)
					w.WriteHeader(204)
					return
				}
				for _, f := range files {
					if f.ID == id {
						if id == "moved-after-list" {
							f.Parents = []string{"other"}
						}
						serveJSON(w, f)
						return
					}
				}
				w.WriteHeader(404)
			}))
			defer srv.Close()
			d := &Drive{client: srv.Client(), base: srv.URL, folder: "folder"}
			n, err := d.Retain(context.Background(), r)
			if err != nil || n != 1 || fmt.Sprint(deleted) != "[expired]" || pages != 2 {
				t.Fatalf("deleted %v count %d pages %d err %v", deleted, n, pages, err)
			}
		})
	}
}
func TestDrivePreflightRejectsNonWritableFolder(t *testing.T) {
	for _, tc := range []struct {
		mime           string
		write, trashed bool
	}{{"application/octet-stream", true, false}, {"application/vnd.google-apps.folder", false, false}, {"application/vnd.google-apps.folder", true, true}} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serveJSON(w, map[string]any{"mimeType": tc.mime, "trashed": tc.trashed, "capabilities": map[string]bool{"canAddChildren": tc.write}})
		}))
		d := &Drive{client: srv.Client(), base: srv.URL, folder: "folder"}
		if d.Preflight(context.Background()) == nil {
			t.Fatal("accepted unwritable destination")
		}
		srv.Close()
	}
}
func TestUploadRangeValidation(t *testing.T) {
	for _, value := range []string{"bytes=2-4", "bytes=0--1", "bytes=0-9223372036854775807", "garbage"} {
		if _, err := uploadedOffset(value); err == nil {
			t.Fatalf("accepted %s", value)
		}
	}
}
