package config

import "testing"

func TestHubUploadLimitsFromEnv(t *testing.T) {
	t.Setenv("HUB_UPLOAD_MAX_BYTES", "2048")
	t.Setenv("HUB_IMAGE_MAX_WIDTH", "800")
	t.Setenv("HUB_IMAGE_MAX_HEIGHT", "600")
	c, err := LoadHubConfig()
	if err != nil || c.UploadMaxBytes != 2048 || c.ImageMaxWidth != 800 || c.ImageMaxHeight != 600 {
		t.Fatalf("Hub upload limits: %+v, %v", c, err)
	}
	t.Setenv("HUB_IMAGE_MAX_WIDTH", "2147483648")
	if _, err := LoadHubConfig(); err == nil {
		t.Fatal("out of range image width accepted")
	}
}
