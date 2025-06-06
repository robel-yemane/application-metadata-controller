package main

import (
	"os"
	"testing"
)

func TestHomeDirUsesHome(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	t.Cleanup(func() {
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserProfile)
	})

	os.Setenv("HOME", "/tmp/home")
	os.Setenv("USERPROFILE", "/tmp/profile")

	got := homeDir()
	if got != "/tmp/home" {
		t.Fatalf("expected %q, got %q", "/tmp/home", got)
	}
}

func TestHomeDirUsesUserProfileWhenHomeEmpty(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	t.Cleanup(func() {
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserProfile)
	})

	os.Setenv("HOME", "")
	os.Setenv("USERPROFILE", "/tmp/profile")

	got := homeDir()
	if got != "/tmp/profile" {
		t.Fatalf("expected %q, got %q", "/tmp/profile", got)
	}
}
