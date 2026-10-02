//go:build integration

package imagesvc

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"go-server/internal/models"
	repository "go-server/internal/repositories"
	"go-server/internal/testutil"
	"go-server/pkg/database"
)

type concurrentStorage struct {
	mu         sync.Mutex
	files      map[string]bool
	repo       repository.ImageRepository
	target     models.ImageTarget
	violations []string
}

func (s *concurrentStorage) Upload(_ context.Context, in models.ImageUploadInput) (*models.StoredImage, error) {
	s.mu.Lock()
	s.files[in.FileName] = true
	s.mu.Unlock()
	return &models.StoredImage{FileID: in.FileName, Path: "/" + in.Folder + "/" + in.FileName, URL: "https://permanent/" + in.FileName}, nil
}
func (s *concurrentStorage) Delete(ctx context.Context, id string) error {
	current, err := s.repo.Read(ctx, s.target)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if current.FileID == id {
		s.violations = append(s.violations, id)
		return fmt.Errorf("attempted deletion of referenced file")
	}
	delete(s.files, id)
	return nil
}
func (s *concurrentStorage) SignedURL(p string, _ time.Time, _ ...models.ImageVariant) (string, error) {
	return "https://signed" + p, nil
}

func TestImagePostgresConcurrencyAndRollback(t *testing.T) {
	ctx := context.Background()
	pg := testutil.StartPostgres(t, ctx)
	t.Cleanup(func() { _ = pg.Terminate(context.Background()) })
	sqlDB, err := sql.Open("pgx", pg.DSN)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	if _, err = migrate.Exec(sqlDB, "postgres", &migrate.FileMigrationSource{Dir: filepath.Join("..", "..", "..", "migrations")}, migrate.Up); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, pg.DSN)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	insert := func(query string, args ...any) int64 {
		t.Helper()
		var id int64
		if err := pool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}
	user := insert(`INSERT INTO users(full_name,email,email_verified) VALUES ('Image Guard','image@example.com',true) RETURNING id`)
	society := insert(`INSERT INTO societies(name,society_code,status,created_by,approved_by,approved_at) VALUES ('Images','IMAGES','active',$1,$1,NOW()) RETURNING id`, user)
	if _, err = pool.Exec(ctx, `INSERT INTO society_members(society_id,user_id,role,status) VALUES ($1,$2,'staff','active')`, society, user); err != nil {
		t.Fatal(err)
	}
	flat := insert(`INSERT INTO flats(society_id,flat_number,created_by) VALUES ($1,'101',$2) RETURNING id`, society, user)
	visitor := insert(`INSERT INTO visitors(full_name,phone_number) VALUES ('Visitor','+919999999999') RETURNING id`)
	entry := insert(`INSERT INTO visitor_entries(society_id,flat_id,visitor_id,source,purpose,status,created_by) VALUES ($1,$2,$3,'guard_entry','guest','approved',$4) RETURNING id`, society, flat, visitor, user)
	database := &database.Database{Pool: pool}
	repo := repository.NewImageRepository(database)
	auth := NewAuthorization(repository.NewUserRepository(database), repository.NewSocietyMemberRepository(database), repository.NewVisitorEntryRepository(database), imageOperational{}, nil)
	for _, target := range []models.ImageTarget{{ActorID: user}, {ActorID: user, SocietyID: society, EntryID: entry, VisitorID: visitor}} {
		kind := "avatar"
		if !target.Avatar() {
			kind = "visitor"
		}
		for _, methods := range [][2]string{{http.MethodPut, http.MethodPut}, {http.MethodPut, http.MethodDelete}, {http.MethodDelete, http.MethodPut}} {
			t.Run(kind+"/"+methods[0]+"-"+methods[1], func(t *testing.T) {
				storage := &concurrentStorage{files: make(map[string]bool), repo: repo, target: target}
				svc := New(storage, repo, auth, "dev/apna-gate", time.Now)
				for range 8 {
					// Begin each round with no reference or provider objects.
					if _, err := svc.Execute(ctx, http.MethodDelete, target, httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/", nil)); err != nil {
						t.Fatal(err)
					}
					var wg sync.WaitGroup
					start := make(chan struct{})
					errs := make(chan error, 2)
					requests := []*http.Request{multipartRequest(t, testImage(t, "png"), false), multipartRequest(t, testImage(t, "jpeg"), false)}
					for i, method := range methods {
						wg.Add(1)
						go func(i int, method string) {
							defer wg.Done()
							<-start
							_, err := svc.Execute(ctx, method, target, httptest.NewRecorder(), requests[i])
							errs <- err
						}(i, method)
					}
					close(start)
					wg.Wait()
					close(errs)
					for err := range errs {
						if err != nil {
							t.Fatal(err)
						}
					}
					current, err := repo.Read(ctx, target)
					if err != nil {
						t.Fatal(err)
					}
					storage.mu.Lock()
					valid := len(storage.violations) == 0 && ((current.FileID == "" && len(storage.files) == 0) || (len(storage.files) == 1 && storage.files[current.FileID]))
					storage.mu.Unlock()
					if !valid {
						t.Fatalf("current=%s files=%v violations=%v", current.FileID, storage.files, storage.violations)
					}
				}
			})
		}
	}
	target := models.ImageTarget{ActorID: user}
	old, err := repo.Read(ctx, target)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.Replace(ctx, target, models.StoredImage{FileID: "rollback", Path: "/rollback", URL: "https://rollback"}, func(context.Context) error { return ErrImageForbidden })
	if err == nil {
		t.Fatal("expected rollback")
	}
	current, err := repo.Read(ctx, target)
	if err != nil || current != old {
		t.Fatal("rollback changed image")
	}
	if _, err = pool.Exec(ctx, `UPDATE users SET avatar_imagekit_file_id='incomplete',avatar_imagekit_file_path=NULL WHERE id=$1`, user); err == nil {
		t.Fatal("incomplete managed reference accepted")
	}

	// Hold an entry lock through a terminal transition; photo replacement must
	// wait, then recheck status and leave the visitor photo untouched.
	visitorTarget := models.ImageTarget{ActorID: user, SocietyID: society, EntryID: entry, VisitorID: visitor}
	before, err := repo.Read(ctx, visitorTarget)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err = tx.Exec(ctx, `UPDATE visitor_entries SET status='checked_out' WHERE id=$1`, entry); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() {
		_, err := repo.Replace(ctx, visitorTarget, models.StoredImage{FileID: "late", Path: "/late", URL: "https://late"}, func(txCtx context.Context) error { _, err := auth.Authorize(txCtx, visitorTarget, true); return err })
		finished <- err
	}()
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if err = <-finished; err == nil {
		t.Fatal("terminal entry accepted replacement")
	}
	after, err := repo.Read(ctx, visitorTarget)
	if err != nil || before != after {
		t.Fatal("state rejection changed photo")
	}
}
