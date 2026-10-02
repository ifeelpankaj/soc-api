package config

import "testing"

func TestImageKitConfig(t *testing.T) {
	valid := ImageKitConfig{Enabled: true, PrivateKey: "test-only", URLEndpoint: "https://ik.imagekit.io/account", FolderPrefix: "prod/apna-gate"}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, prefix := range []string{"", "/prod", "prod/", "prod/../a", "prod//a", "prod\\a", "prod/%2e"} {
		cfg := valid
		cfg.FolderPrefix = prefix
		if cfg.Validate() == nil {
			t.Fatalf("accepted %q", prefix)
		}
	}
	for _, endpoint := range []string{"", "http://ik.imagekit.io/account", "https://user:pass@host/path", "https://host/path?key=x", "https://host/#x"} {
		cfg := valid
		cfg.URLEndpoint = endpoint
		if cfg.Validate() == nil {
			t.Fatalf("accepted %q", endpoint)
		}
	}
	valid.PrivateKey = ""
	if valid.Validate() == nil {
		t.Fatal("accepted missing key")
	}
	if err := (ImageKitConfig{}).Validate(); err != nil {
		t.Fatal(err)
	}
}
