package config

import (
	"testing"
	"time"
)

func TestDriveFolderParsing(t *testing.T) {
	for _, value := range []string{"folder_123-abc", "https://drive.google.com/drive/folders/folder_123-abc?usp=sharing", "https://drive.google.com/drive/u/0/folders/folder_123-abc"} {
		id, err := DriveFolderID(value)
		if err != nil || id != "folder_123-abc" {
			t.Fatalf("%q -> %q %v", value, id, err)
		}
	}
	for _, value := range []string{"", "https://evil.com/drive/folders/id", "https://drive.google.com.evil.com/folders/id", "http://drive.google.com/folders/id", "https://user@drive.google.com/folders/id", "https://drive.google.com/folders/id/other", "../id"} {
		if _, err := DriveFolderID(value); err == nil {
			t.Fatalf("accepted %q", value)
		}
	}
}
func TestBackupConfigRequiresCredentialsOnlyWhenEnabled(t *testing.T) {
	if err := (BackupConfig{}).Validate(); err != nil {
		t.Fatal(err)
	}
	c := BackupConfig{Enabled: true, Deployment: "production", TempDir: t.TempDir(), Timeout: 2 * time.Hour, Folder: "folder", ClientID: "client", ClientSecret: "secret", RefreshToken: "refresh"}
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*BackupConfig){func(c *BackupConfig) { c.RefreshToken = "" }, func(c *BackupConfig) { c.TempDir = "relative" }, func(c *BackupConfig) { c.Timeout = 0 }, func(c *BackupConfig) { c.Deployment = "" }} {
		bad := c
		mutate(&bad)
		if bad.Validate() == nil {
			t.Fatal("accepted invalid configuration")
		}
	}
}
