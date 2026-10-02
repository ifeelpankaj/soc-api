package imagesvc

import (
	"context"
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec // ImageKit's URL-signing protocol requires HMAC-SHA1.
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	imagekit "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"go-server/internal/config"
	"go-server/internal/models"
)

func TestImageKitVariantSignatures(t *testing.T) {
	s := NewImageKitStorage(config.ImageKitConfig{PrivateKey: "test-only", URLEndpoint: "https://ik.imagekit.io/account", FolderPrefix: "dev/apna-gate"})
	expires := time.Now().Add(15 * time.Minute).Truncate(time.Second)
	for variant, transform := range map[models.ImageVariant]string{
		models.ImageOriginal: "", models.ImageList: "tr:w-80,h-80,fo-face/", models.ImageAvatar: "tr:w-300,h-300,fo-face/", models.ImageDetail: "tr:w-600/",
	} {
		signed, err := s.SignedURL("/dev/apna-gate/photo.jpg", expires, variant)
		if err != nil {
			t.Fatal(err)
		}
		u, err := url.Parse(signed)
		if err != nil {
			t.Fatal(err)
		}
		path := strings.TrimPrefix(u.Path, "/account/")
		if path != transform+"dev/apna-gate/photo.jpg" {
			t.Fatalf("%s: %s", variant, path)
		}
		mac := hmac.New(sha1.New, []byte("test-only"))
		mac.Write([]byte(path + strconv.FormatInt(expires.Unix(), 10)))
		if u.Query().Get("ik-s") != hex.EncodeToString(mac.Sum(nil)) {
			t.Fatal("transformation not included in signature")
		}
	}
	if _, err := s.SignedURL("/dev/apna-gate/photo.jpg", expires, "w-999"); !errors.Is(err, ErrImageVariant) {
		t.Fatal(err)
	}
	for _, value := range []string{"", "original", "list", "avatar", "detail"} {
		if _, err := ParseImageVariant(value); err != nil {
			t.Fatal(err)
		}
	}
}

type imageRoundTripper func(*http.Request) (*http.Response, error)

func (f imageRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestImageKitAdapterPrivateUploadAndNoRetries(t *testing.T) {
	for _, status := range []int{200, 500} {
		calls := 0
		transport := imageRoundTripper(func(req *http.Request) (*http.Response, error) {
			calls++
			req.Body = http.MaxBytesReader(nil, req.Body, models.ImageMaxBodyBytes)
			reader, err := req.MultipartReader()
			if err != nil {
				t.Fatal(err)
			}
			fields := map[string]string{}
			for {
				part, err := reader.NextPart()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				data, err := io.ReadAll(io.LimitReader(part, models.ImageMaxBodyBytes+1))
				if err != nil {
					t.Fatal(err)
				}
				fields[part.FormName()] = string(data)
			}
			if fields["isPrivateFile"] != "true" || fields["overwriteFile"] != "false" || fields["useUniqueFileName"] != "false" {
				t.Fatal("unsafe upload settings")
			}
			if fields["fileName"] != "generated.png" || fields["folder"] != "dev/apna-gate/users/1/avatars" {
				t.Fatal("wrong location")
			}
			body := `{"fileId":"new","filePath":"/dev/apna-gate/users/1/avatars/generated.png","url":"https://ik.imagekit.io/account/dev/apna-gate/users/1/avatars/generated.png","isPrivateFile":true}`
			if status == 500 {
				body = `{"message":"private-provider-details"}`
			}
			return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
		})
		storage := NewImageKitStorage(config.ImageKitConfig{PrivateKey: "test-only", URLEndpoint: "https://ik.imagekit.io/account", FolderPrefix: "dev/apna-gate"})
		storage.client = imagekit.NewClient(option.WithPrivateKey("test-only"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: transport}))
		result, err := storage.Upload(context.Background(), models.ImageUploadInput{Reader: strings.NewReader("bytes"), FileName: "generated.png", Folder: "dev/apna-gate/users/1/avatars"})
		if calls != 1 {
			t.Fatalf("retried %d times", calls)
		}
		if status == 200 && (err != nil || result.FileID != "new") {
			t.Fatalf("%v %v", result, err)
		}
		if status == 500 && (err == nil || strings.Contains(err.Error(), "private-provider")) {
			t.Fatalf("unsanitized error %v", err)
		}
	}
}

func TestImageKitSigningAndDelete(t *testing.T) {
	storage := NewImageKitStorage(config.ImageKitConfig{PrivateKey: "test-only", URLEndpoint: "https://ik.imagekit.io/account", FolderPrefix: "dev/apna-gate"})
	expiry := time.Now().Add(models.ImageSignedURLTTL).Truncate(time.Second)
	signed, err := storage.SignedURL("/dev/apna-gate/users/1/avatars/a.png", expiry)
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(signed)
	if err != nil {
		t.Fatal(err)
	}
	if u.Query().Get("ik-t") != strconv.FormatInt(expiry.Unix(), 10) || u.Query().Get("ik-s") == "" {
		t.Fatal("missing signature/expiry")
	}
	for _, p := range []string{"https://evil/image", "/other/image", "/dev/apna-gate/../image", "/dev/apna-gate/%2e%2e/image", "/dev/apna-gate/a?ik-t=999"} {
		if _, err := storage.SignedURL(p, expiry); err == nil {
			t.Fatalf("signed %s", p)
		}
	}
	storage.client = imagekit.NewClient(option.WithPrivateKey("test-only"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: imageRoundTripper(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 404, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"message":"absent"}`)), Request: req}, nil
	})}))
	if err := storage.Delete(context.Background(), "absent"); err != nil {
		t.Fatal(err)
	}
}
