package imagesvc

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-server/internal/models"
	repository "go-server/internal/repositories"
)

type memoryImageRepo struct {
	stored             models.StoredImage
	err, resolveErr    error
	commitDespiteError bool
	locked             bool
	calls              int
}

func (r *memoryImageRepo) Read(context.Context, models.ImageTarget) (models.StoredImage, error) {
	return r.stored, nil
}
func (r *memoryImageRepo) Resolve(context.Context, models.ImageTarget) (models.StoredImage, error) {
	return r.stored, r.resolveErr
}
func (r *memoryImageRepo) Replace(ctx context.Context, _ models.ImageTarget, next models.StoredImage, auth func(context.Context) error) (models.StoredImage, error) {
	r.calls++
	r.locked = true
	defer func() { r.locked = false }()
	old := r.stored
	if err := auth(ctx); err != nil {
		return old, err
	}
	if r.err == nil || r.commitDespiteError {
		r.stored = next
	}
	return old, r.err
}

type fakeImageAuth struct {
	calls       int
	deniedAfter int
	err         error
}

func (a *fakeImageAuth) Authorize(_ context.Context, t models.ImageTarget, _ bool) (models.ImageTarget, error) {
	a.calls++
	if a.err != nil && a.calls > a.deniedAfter {
		return t, a.err
	}
	return t, nil
}

type fakeImageStorage struct {
	t                    *testing.T
	repo                 *memoryImageRepo
	uploadErr, deleteErr error
	deleted              []string
	onUpload             func()
	expires              time.Time
	uploaded             int
	input                models.ImageUploadInput
	data                 []byte
}

func (s *fakeImageStorage) Upload(ctx context.Context, in models.ImageUploadInput) (*models.StoredImage, error) {
	s.t.Helper()
	if s.repo.locked {
		s.t.Fatal("upload under database lock")
	}
	s.uploaded++
	s.input = in
	var err error
	s.data, err = io.ReadAll(in.Reader)
	if err != nil {
		s.t.Fatal(err)
	}
	if s.onUpload != nil {
		s.onUpload()
	}
	return &models.StoredImage{FileID: "new", Path: "/dev/apna-gate/new.png", URL: "https://permanent/new"}, s.uploadErr
}

func TestVisitorNormalizedButAvatarPreserved(t *testing.T) {
	for _, entryID := range []int64{0, 3} {
		repo := &memoryImageRepo{}
		storage := &fakeImageStorage{t: t, repo: repo}
		svc := New(storage, repo, &fakeImageAuth{}, "dev/apna-gate", nil)
		input := testImage(t, "png")
		_, err := svc.Execute(context.Background(), http.MethodPut, models.ImageTarget{ActorID: 1, SocietyID: 2, EntryID: entryID}, httptest.NewRecorder(), multipartRequest(t, input, false))
		if err != nil {
			t.Fatal(err)
		}
		if entryID == 0 {
			if !bytes.Equal(input, storage.data) {
				t.Fatal("avatar bytes changed")
			}
		} else {
			_, format, err := image.DecodeConfig(bytes.NewReader(storage.data))
			if err != nil || format != "jpeg" || !strings.HasSuffix(storage.input.FileName, ".jpg") {
				t.Fatalf("visitor upload was not normalized: %s %v", format, err)
			}
		}
	}
}

func TestViewingVariantContract(t *testing.T) {
	repo := &memoryImageRepo{stored: models.StoredImage{FileID: "id", Path: "/path", URL: "private"}}
	storage := &fakeImageStorage{t: t, repo: repo}
	svc := New(storage, repo, &fakeImageAuth{}, "dev/apna-gate", nil)
	for _, variant := range []string{"", "list", "avatar", "detail", "original", "w-9999"} {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/?variant="+variant, nil)
		view, err := svc.Execute(context.Background(), http.MethodGet, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), req)
		if variant == "w-9999" {
			if !errors.Is(err, ErrImageVariant) {
				t.Fatal(err)
			}
			continue
		}
		want := variant
		if want == "" {
			want = "original"
		}
		if err != nil || string(view.Variant) != want {
			t.Fatalf("%s: %v %v", variant, view, err)
		}
	}
}
func (s *fakeImageStorage) Delete(ctx context.Context, id string) error {
	s.t.Helper()
	if s.repo.locked {
		s.t.Fatal("delete under database lock")
	}
	if ctx.Err() != nil {
		s.t.Fatal("cleanup inherited canceled request")
	}
	if id == s.repo.stored.FileID {
		s.t.Fatal("deleted currently referenced file")
	}
	s.deleted = append(s.deleted, id)
	return s.deleteErr
}
func (s *fakeImageStorage) SignedURL(_ string, expires time.Time, _ ...models.ImageVariant) (string, error) {
	s.expires = expires
	return "https://signed/view", nil
}

func TestImageReplacementFailures(t *testing.T) {
	for _, tc := range []struct {
		name                string
		dbErr, errorCleanup error
		commit              bool
		resolveErr          error
		wantCurrent         string
		wantError           bool
		wantDeleted         string
	}{
		{name: "success", wantCurrent: "new", wantDeleted: "old"},
		{name: "cleanup failure does not fail commit", errorCleanup: errors.New("secret provider detail"), wantCurrent: "new", wantDeleted: "old"},
		{name: "rollback cleanup", dbErr: errors.New("db secret"), wantCurrent: "old", wantError: true, wantDeleted: "new"},
		{name: "rollback cleanup failure", dbErr: errors.New("db"), errorCleanup: errors.New("provider"), wantCurrent: "old", wantError: true, wantDeleted: "new"},
		{name: "uncertain committed", dbErr: &repository.ImageCommitUncertain{}, commit: true, wantCurrent: "new", wantDeleted: "old"},
		{name: "uncertain rolled back", dbErr: &repository.ImageCommitUncertain{}, wantCurrent: "old", wantError: true, wantDeleted: "new"},
		{name: "unresolved commit retains both", dbErr: &repository.ImageCommitUncertain{}, commit: true, resolveErr: errors.New("offline"), wantCurrent: "new", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &memoryImageRepo{stored: models.StoredImage{FileID: "old", Path: "/old", URL: "old"}, err: tc.dbErr, commitDespiteError: tc.commit, resolveErr: tc.resolveErr}
			storage := &fakeImageStorage{t: t, repo: repo, deleteErr: tc.errorCleanup}
			now := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
			svc := New(storage, repo, &fakeImageAuth{}, "dev/apna-gate", func() time.Time { return now })
			view, err := svc.Execute(context.Background(), http.MethodPut, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), multipartRequest(t, testImage(t, "png"), false))
			if (err != nil) != tc.wantError {
				t.Fatalf("err=%v", err)
			}
			if repo.stored.FileID != tc.wantCurrent {
				t.Fatalf("current %s", repo.stored.FileID)
			}
			if strings.Join(storage.deleted, ",") != tc.wantDeleted {
				t.Fatalf("deleted %v", storage.deleted)
			}
			if err == nil && (view.URL != "https://signed/view" || !view.ExpiresAt.Equal(now.Add(models.ImageSignedURLTTL))) {
				t.Fatalf("view %#v", view)
			}
			if err != nil && strings.Contains(err.Error(), "secret") {
				t.Fatal("leaked detail")
			}
		})
	}
}

func TestAuthorizationBeforeReadingAndAfterUpload(t *testing.T) {
	for _, after := range []int{0, 1} {
		repo := &memoryImageRepo{}
		storage := &fakeImageStorage{t: t, repo: repo}
		auth := &fakeImageAuth{deniedAfter: after, err: ErrImageState}
		svc := New(storage, repo, auth, "dev/apna-gate", nil)
		req := multipartRequest(t, testImage(t, "png"), false)
		_, err := svc.Execute(context.Background(), http.MethodPut, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), req)
		if !errors.Is(err, ErrImageState) {
			t.Fatal(err)
		}
		if storage.uploaded != after {
			t.Fatalf("uploads %d", storage.uploaded)
		}
		if repo.stored.Managed() {
			t.Fatal("saved despite authorization failure")
		}
	}
}

func TestLegacyAndDeletion(t *testing.T) {
	repo := &memoryImageRepo{stored: models.StoredImage{URL: "https://legacy"}}
	storage := &fakeImageStorage{t: t, repo: repo}
	svc := New(storage, repo, &fakeImageAuth{}, "dev/apna-gate", nil)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	if _, err := svc.Execute(context.Background(), http.MethodGet, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), req); !errors.Is(err, ErrImageNotFound) {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := svc.Execute(context.Background(), http.MethodDelete, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), req); err != nil {
			t.Fatal(err)
		}
	}
	if repo.stored.URL != "" || len(storage.deleted) != 0 {
		t.Fatal("legacy deletion contacted provider")
	}
}

func TestCancellationCompensatesWithIndependentContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repo := &memoryImageRepo{err: context.Canceled}
	storage := &fakeImageStorage{t: t, repo: repo, onUpload: cancel}
	svc := New(storage, repo, &fakeImageAuth{}, "dev/apna-gate", nil)
	_, err := svc.Execute(ctx, http.MethodPut, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), multipartRequest(t, testImage(t, "png"), false))
	if err == nil || len(storage.deleted) != 1 {
		t.Fatalf("err %v deleted %v", err, storage.deleted)
	}
}

func TestDisabledAndBusy(t *testing.T) {
	repo := &memoryImageRepo{}
	auth := &fakeImageAuth{}
	svc := New(nil, repo, auth, "dev/apna-gate", nil)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/", nil)
	if _, err := svc.Execute(context.Background(), http.MethodPut, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), req); !errors.Is(err, ErrImageUnavailable) {
		t.Fatal(err)
	}
	svc.storage = &fakeImageStorage{t: t, repo: repo}
	for range models.ImageConcurrency {
		svc.slots <- struct{}{}
	}
	if _, err := svc.Execute(context.Background(), http.MethodPut, models.ImageTarget{ActorID: 1}, httptest.NewRecorder(), req); !errors.Is(err, ErrImageBusy) {
		t.Fatal(err)
	}
}
