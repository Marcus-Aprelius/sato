package sato_test

import (
	"os"
	"path/filepath"
	"testing"

	"sato/internal/sato"
)

// TestFindKDBX ensures that FindKDBX picks up the first *.kdbx file in a
// directory and returns "" for empty or missing directories.
func TestFindKDBX(t *testing.T) {
	dir := t.TempDir()

	if got := sato.FindKDBX(dir); got != "" {
		t.Fatalf("empty dir: got %q; want empty string", got)
	}

	target := filepath.Join(dir, "secrets.kdbx")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if got := sato.FindKDBX(dir); got != target {
		t.Fatalf("populated dir: got %q; want %q", got, target)
	}

	missing := filepath.Join(dir, "does-not-exist")
	if got := sato.FindKDBX(missing); got != "" {
		t.Fatalf("missing dir: got %q; want empty string", got)
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
