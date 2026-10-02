package config

import (
	"fmt"
	"os"
	"strconv"
)

type HubConfig struct {
	PostsPerWindow, PostWindowSeconds, CommentsPerWindow, CommentWindowSeconds int32
	UploadMaxBytes                                                             int64
	ImageMaxWidth, ImageMaxHeight                                              int32
}

func LoadHubConfig() (HubConfig, error) {
	c := HubConfig{PostsPerWindow: 5, PostWindowSeconds: 600, CommentsPerWindow: 30, CommentWindowSeconds: 60, UploadMaxBytes: 10485760, ImageMaxWidth: 1600, ImageMaxHeight: 1600}
	for key, dst := range map[string]*int32{"HUB_POST_LIMIT": &c.PostsPerWindow, "HUB_POST_WINDOW_SECONDS": &c.PostWindowSeconds, "HUB_COMMENT_LIMIT": &c.CommentsPerWindow, "HUB_COMMENT_WINDOW_SECONDS": &c.CommentWindowSeconds} {
		if value, ok := os.LookupEnv(key); ok {
			n, err := strconv.ParseInt(value, 10, 32)
			if err != nil || n < 1 || n > 86400 {
				return c, fmt.Errorf("%s must be between 1 and 86400", key)
			}
			*dst = int32(n)
		}
	}
	if value, ok := os.LookupEnv("HUB_UPLOAD_MAX_BYTES"); ok {
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil || n < 1 || n > 100*1024*1024 {
			return c, fmt.Errorf("HUB_UPLOAD_MAX_BYTES must be between 1 and 104857600")
		}
		c.UploadMaxBytes = n
	}
	for key, dst := range map[string]*int32{"HUB_IMAGE_MAX_WIDTH": &c.ImageMaxWidth, "HUB_IMAGE_MAX_HEIGHT": &c.ImageMaxHeight} {
		if value, ok := os.LookupEnv(key); ok {
			n, err := strconv.ParseInt(value, 10, 32)
			if err != nil || n < 1 || n > 32767 {
				return c, fmt.Errorf("%s must be between 1 and 32767", key)
			}
			*dst = int32(n)
		}
	}
	return c, nil
}
