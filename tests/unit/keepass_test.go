package sato_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"sato/internal/sato"
	"sato/internal/testdb"
)

const testDBPassword = testdb.DefaultPassword

func projectRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get wd: %v", err)
	}

	root, err := filepath.Abs(filepath.Join(wd, "..", ".."))
	if err != nil {
		t.Fatalf("project root: %v", err)
	}

	return root
}

func createUnitKeePassDB(t *testing.T) string {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "unit-secrets.kdbx")

	if err := testdb.WriteDemoDatabase(dbPath, testDBPassword); err != nil {
		t.Fatalf("create test db: %v", err)
	}

	return dbPath
}

func TestLoadKeePass_Playground(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	entries, err := sato.LoadKeePass(dbPath, testDBPassword)
	if err != nil {
		t.Fatalf("LoadKeePass: %v", err)
	}

	want := map[string]string{
		"DB_PASSWORD": "super_secret_db_password_123",
		"API_KEY":     "sk-test-api-key-xyz789",
	}

	for k, v := range want {
		got, ok := entries[k]
		if !ok {
			t.Fatalf("missing entry %q", k)
		}

		if got != v {
			t.Fatalf("entry %q: got %q; want %q", k, got, v)
		}
	}
}

func TestLoadKeePassEntries_ReadsNestedGroups(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	entries, err := sato.LoadKeePassEntries(dbPath, testDBPassword)
	if err != nil {
		t.Fatalf("LoadKeePassEntries: %v", err)
	}

	got := make(map[string]string)

	for _, entry := range entries {
		got[entry.Path] = entry.Value
	}

	want := map[string]string{
		"Secrets/API_KEY":        "sk-test-api-key-xyz789",
		"Secrets/DB_PASSWORD":    "super_secret_db_password_123",
		"Secrets/NGINX_PASSWORD": "nginx_password_demo_456",
		"Secrets/test/test1":     "test_secret_value_1",
		"Secrets/test/test2":     "test_secret_value_2",
	}

	for path, value := range want {
		gotValue, ok := got[path]
		if !ok {
			t.Fatalf("missing path %q; got paths: %v", path, sortedKeys(got))
		}

		if gotValue != value {
			t.Fatalf("path %q: got %q; want %q", path, gotValue, value)
		}
	}
}

func TestLoadKeePassData_DetectsEmptyGroups(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	data, err := sato.LoadKeePassData(dbPath, testDBPassword)
	if err != nil {
		t.Fatalf("LoadKeePassData: %v", err)
	}

	got := make(map[string]bool)

	for _, path := range data.EmptyGroupPaths {
		got[path] = true
	}

	want := []string{
		"Secrets/empty-group",
		"Secrets/test/empty-nested",
	}

	for _, path := range want {
		if !got[path] {
			t.Fatalf("missing empty group %q; got %v", path, data.EmptyGroupPaths)
		}
	}
}

func TestLoadKeePassData_DoesNotMarkGroupsWithSecretsAsEmpty(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	data, err := sato.LoadKeePassData(dbPath, testDBPassword)
	if err != nil {
		t.Fatalf("LoadKeePassData: %v", err)
	}

	notEmpty := map[string]bool{
		"Secrets":      true,
		"Secrets/test": true,
	}

	for _, path := range data.EmptyGroupPaths {
		if notEmpty[path] {
			t.Fatalf("group %q must not be reported as empty", path)
		}
	}
}

func TestLoadKeePass_MapKeepsSecretNames(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	env, err := sato.LoadKeePass(dbPath, testDBPassword)
	if err != nil {
		t.Fatalf("LoadKeePass: %v", err)
	}

	want := map[string]string{
		"DB_PASSWORD":    "super_secret_db_password_123",
		"API_KEY":        "sk-test-api-key-xyz789",
		"NGINX_PASSWORD": "nginx_password_demo_456",
		"test1":          "test_secret_value_1",
		"test2":          "test_secret_value_2",
	}

	for key, value := range want {
		got, ok := env[key]
		if !ok {
			t.Fatalf("missing env key %q", key)
		}

		if got != value {
			t.Fatalf("env key %q: got %q; want %q", key, got, value)
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
