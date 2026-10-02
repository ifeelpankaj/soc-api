package config

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type BackupConfig struct {
	Enabled      bool
	Deployment   string
	TempDir      string
	Timeout      time.Duration
	Folder       string
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func loadBackupConfig() BackupConfig {
	return BackupConfig{
		Enabled:      getEnvAsBool("BACKUP_ENABLED", false),
		Deployment:   strings.TrimSpace(getEnv("BACKUP_DEPLOYMENT", "")),
		TempDir:      getEnv("BACKUP_TEMP_DIR", "/var/lib/apna-gate/backups/tmp"),
		Timeout:      time.Duration(getEnvAsInt("BACKUP_TIMEOUT_SECONDS", 7200)) * time.Second,
		Folder:       strings.TrimSpace(getEnv("BACKUP_GDRIVE_FOLDER", "")),
		ClientID:     strings.TrimSpace(getEnv("BACKUP_GDRIVE_CLIENT_ID", "")),
		ClientSecret: strings.TrimSpace(getEnv("BACKUP_GDRIVE_CLIENT_SECRET", "")),
		RefreshToken: strings.TrimSpace(getEnv("BACKUP_GDRIVE_REFRESH_TOKEN", "")),
	}
}

var driveID = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// DriveFolderID parses an ID or a Google Drive folders URL, never an arbitrary URL.
func DriveFolderID(value string) (string, error) {
	value = strings.TrimSpace(value)
	if driveID.MatchString(value) {
		return value, nil
	}
	u, err := url.Parse(value)
	if err == nil && u.Scheme == "https" && u.Host == "drive.google.com" && u.User == nil {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		for i, part := range parts {
			if part == "folders" && i == len(parts)-2 && driveID.MatchString(parts[i+1]) {
				return parts[i+1], nil
			}
		}
	}
	return "", fmt.Errorf("BACKUP_GDRIVE_FOLDER must be a Drive folder ID or https://drive.google.com/.../folders/ID URL")
}

func (c BackupConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Deployment == "" || len(c.Deployment) > 60 {
		return fmt.Errorf("BACKUP_DEPLOYMENT must contain 1-60 characters")
	}
	if !filepath.IsAbs(c.TempDir) || filepath.Dir(filepath.Clean(c.TempDir)) == filepath.Clean(c.TempDir) {
		return fmt.Errorf("BACKUP_TEMP_DIR must be an absolute dedicated directory")
	}
	if c.Timeout < time.Minute || c.Timeout > 24*time.Hour {
		return fmt.Errorf("BACKUP_TIMEOUT_SECONDS must be between 60 and 86400")
	}
	if _, err := DriveFolderID(c.Folder); err != nil {
		return err
	}
	if c.ClientID == "" || c.ClientSecret == "" || c.RefreshToken == "" {
		return fmt.Errorf("backup Google Drive OAuth credentials are required when BACKUP_ENABLED=true")
	}
	return nil
}
