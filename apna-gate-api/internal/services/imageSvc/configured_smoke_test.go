package imagesvc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go-server/internal/config"
	"go-server/internal/models"
)

// Opt-in only: creates a synthetic private file and deletes it without touching
// visitor/profile records. Never print provider errors, credentials or URLs.
func TestConfiguredImageKitSmoke(t *testing.T) {
	if os.Getenv("APNA_GATE_IMAGEKIT_SMOKE") != "1" {
		t.Skip("requires explicit ImageKit smoke-test opt-in")
	}
	values, err := godotenv.Read("../../../.env.development")
	if err != nil {
		t.Fatal("development configuration unavailable")
	}
	if values["IMAGEKIT_ENABLED"] != "true" {
		t.Fatal("ImageKit disabled in development configuration")
	}
	storage := NewImageKitStorage(config.ImageKitConfig{Enabled: true, PrivateKey: values["IMAGEKIT_PRIVATE_KEY"], URLEndpoint: values["IMAGEKIT_URL_ENDPOINT"], FolderPrefix: values["IMAGEKIT_FOLDER_PREFIX"]})
	data, err := normalizeVisitorImage(testImage(t, "png"))
	if err != nil {
		t.Fatal("normalization failed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stored, err := storage.Upload(ctx, models.ImageUploadInput{Reader: bytes.NewReader(data), FileName: uuid.NewString() + ".jpg", Folder: values["IMAGEKIT_FOLDER_PREFIX"] + "/verification"})
	if stored != nil && stored.FileID != "" {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			if storage.Delete(ctx, stored.FileID) != nil {
				t.Error("test file cleanup failed")
			}
		}()
	}
	if err != nil {
		t.Fatal("configured ImageKit upload failed")
	}
	for _, variant := range []models.ImageVariant{models.ImageOriginal, models.ImageList, models.ImageAvatar, models.ImageDetail} {
		signed, err := storage.SignedURL(stored.Path, time.Now().Add(15*time.Minute), variant)
		if err != nil {
			t.Fatalf("signing %s failed", variant)
		}
		checkDelivery := func(address string, allowed bool) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
			if err != nil {
				t.Fatal("delivery request could not be constructed")
			}
			response, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal("delivery request failed")
			}
			_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, models.ImageMaxUploadBytes))
			if err = response.Body.Close(); err != nil {
				t.Fatal("delivery response body could not be closed")
			}
			if allowed && response.StatusCode != 200 {
				t.Fatalf("%s delivery returned %d", variant, response.StatusCode)
			}
			if !allowed && response.StatusCode != 401 && response.StatusCode != 403 {
				t.Fatalf("private access rejection returned %d", response.StatusCode)
			}
		}
		checkDelivery(signed, true)
		parsed, _ := url.Parse(signed)
		query := parsed.Query()
		query.Set("ik-s", "invalid-signature")
		parsed.RawQuery = query.Encode()
		checkDelivery(parsed.String(), false)
		parsed.RawQuery = ""
		checkDelivery(parsed.String(), false)
	}
}
