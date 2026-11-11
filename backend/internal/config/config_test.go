package config

import "testing"

func TestFromEnvironmentMissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	_, err := FromEnvironment()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is empty")
	}
}

func TestFromEnvironmentDefaultPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("PORT", "")
	got, err := FromEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if got.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", got.Port)
	}
}
