package hubsvc

import (
	"bytes"
	"go-server/internal/config"
	"image"
	"image/png"
	"strings"
	"testing"
	"time"
)

func TestCursorAndValidation(t *testing.T) {
	now := time.Now().UTC()
	raw := encodeCursor("society/channel", now, 45)
	c, e := decodeCursor(raw, "society/channel")
	if e != nil || c.ID != 45 || !c.Time.Equal(now) {
		t.Fatalf("round trip: %+v %v", c, e)
	}
	for _, v := range []string{"bad", raw, strings.Repeat("a", 1025)} {
		if _, e := decodeCursor(v, "other/channel"); e == nil {
			t.Fatal("accepted invalid/cross-resource cursor")
		}
	}
	for _, n := range []int{-1, 101} {
		if _, e := pageSize(n); e == nil {
			t.Fatal("accepted invalid page size")
		}
	}
	if validAttachments(ptr([]int64{1, 1})) || validAttachments(ptr([]int64{0})) || validBody("  ", 20) {
		t.Fatal("invalid content accepted")
	}
	if !validBody("नमस्ते", 20) {
		t.Fatal("unicode rejected")
	}
}
func TestMediaValidation(t *testing.T) {
	limits := config.HubConfig{UploadMaxBytes: 10485760, ImageMaxWidth: 1600, ImageMaxHeight: 1600}
	var b bytes.Buffer
	if e := png.Encode(&b, image.NewRGBA(image.Rect(0, 0, 2, 3))); e != nil {
		t.Fatal(e)
	}
	mime, ext, w, h, e := validateFile(b.Bytes(), limits)
	if e != nil || mime != "image/png" || ext != ".png" || *w != 2 || *h != 3 {
		t.Fatalf("PNG validation: %s %s %v", mime, ext, e)
	}
	tooNarrow := limits
	tooNarrow.ImageMaxWidth = 1
	if _, _, _, _, e := validateFile(b.Bytes(), tooNarrow); e == nil {
		t.Fatal("image wider than configured limit accepted")
	}
	tooSmall := limits
	tooSmall.UploadMaxBytes = int64(b.Len() - 1)
	if _, _, _, _, e := validateFile(b.Bytes(), tooSmall); e == nil {
		t.Fatal("file larger than configured limit accepted")
	}
	mime, ext, w, h, e = validateFile([]byte("%PDF-1.7\n1 0 obj\n<<>>\nendobj\n%%EOF"), limits)
	if e != nil || mime != "application/pdf" || ext != ".pdf" || w != nil || h != nil {
		t.Fatal("PDF validation failed", e)
	}
	for _, data := range [][]byte{nil, []byte("<svg></svg>"), []byte("%PDF-1.7\ntruncated"), b.Bytes()[:20], make([]byte, 10*1024*1024+1)} {
		if _, _, _, _, e := validateFile(data, limits); e == nil {
			t.Fatal("invalid media accepted")
		}
	}
}
