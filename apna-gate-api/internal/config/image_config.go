package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type ImageKitConfig struct {
	Enabled      bool
	PrivateKey   string
	URLEndpoint  string
	FolderPrefix string
}

func (c ImageKitConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if strings.TrimSpace(c.PrivateKey) == "" {
		return fmt.Errorf("IMAGEKIT_PRIVATE_KEY is required")
	}
	u, err := url.Parse(c.URLEndpoint)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("IMAGEKIT_URL_ENDPOINT must be an HTTPS URL without credentials, query, or fragment")
	}
	if c.FolderPrefix == "" || strings.HasPrefix(c.FolderPrefix, "/") || strings.HasSuffix(c.FolderPrefix, "/") {
		return fmt.Errorf("IMAGEKIT_FOLDER_PREFIX must be a relative folder prefix")
	}
	for _, part := range strings.Split(c.FolderPrefix, "/") {
		if part == "" || part == "." || part == ".." {
			return fmt.Errorf("invalid IMAGEKIT_FOLDER_PREFIX")
		}
		for _, ch := range part {
			allowed := ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-' || ch == '_'
			if !allowed {
				return fmt.Errorf("invalid IMAGEKIT_FOLDER_PREFIX")
			}
		}
	}
	return nil
}

func (c *Config) validateImages() error {
	if err := c.ImageKit.Validate(); err != nil {
		return err
	}
	if c.ImageKit.Enabled && (c.ReadTimeout < 60*time.Second || c.WriteTimeout < 90*time.Second) {
		return fmt.Errorf("ImageKit requires READ_TIMEOUT >= 60 and WRITE_TIMEOUT >= 90 seconds")
	}
	return nil
}
