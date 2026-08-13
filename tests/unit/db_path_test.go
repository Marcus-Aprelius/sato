package sato_test

import (
	"os"
	"path/filepath"
	"testing"

	"sato/internal/sato"
)

// TestFindSecretDB ensures that FindSecretDB picks up the first *.kdbx file in a
// directory and returns "" for empty or missing directories.
func TestFindSecretDB(t *testing.T) {
	tests := []struct {
		name      string
		extension string
	}{
		{"kdbx", ".kdbx"},
		{"psafe3", ".psafe3"},
		{"ibak", ".ibak"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, "secrets"+tt.extension)

			if err := os.WriteFile(target, []byte("test"), 0o600); err != nil {
				t.Fatalf("write fixture: %v", err)
			}

			if got := sato.FindSecretDB(dir); got != target {
				t.Fatalf("FindSecretDB() = %q; want %q", got, target)
			}
		})
	}
}

// TestFindDBPath_Priority verifies the priority order:
// 1) --db-path flag, 2) ~/.sato/*.kdbx, 3) SATO_DB_PATH env var.
func TestFindDBPath_Priority(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	flagPath := filepath.Join(t.TempDir(), "flag.kdbx")
	envPath := filepath.Join(t.TempDir(), "env.kdbx")
	localDir := filepath.Join(home, ".sato")
	localPath := filepath.Join(localDir, "local.kdbx")

	if err := os.MkdirAll(localDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	for _, p := range []string{flagPath, envPath, localPath} {
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatalf("write %s: %v", p, err)
		}
	}

	t.Setenv("SATO_DB_PATH", envPath)

	if got := sato.FindDBPath(flagPath); got != flagPath {
		t.Fatalf("priority 1: got %q; want %q", got, flagPath)
	}

	if got := sato.FindDBPath(""); got != localPath {
		t.Fatalf("priority 2: got %q; want %q", got, localPath)
	}

	os.RemoveAll(localDir)

	if got := sato.FindDBPath(""); got != envPath {
		t.Fatalf("priority 3: got %q; want %q", got, envPath)
	}

	t.Setenv("SATO_DB_PATH", "")

	if got := sato.FindDBPath(""); got != "" {
		t.Fatalf("nothing set: got %q; want empty string", got)
	}
}

func TestFindDBPath_IgnoresMissingPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	missingFlag := filepath.Join(t.TempDir(), "missing-flag.kdbx")
	missingEnv := filepath.Join(t.TempDir(), "missing-env.kdbx")

	t.Setenv("SATO_DB_PATH", missingEnv)

	if got := sato.FindDBPath(missingFlag); got != "" {
		t.Fatalf("missing paths: got %q; want empty string", got)
	}
}
