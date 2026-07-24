package sato_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"sato/internal/sato"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

const testDBPassword = "sato"

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

func playgroundDBPath(t *testing.T) string {
	t.Helper()

	path := filepath.Join(projectRoot(t), "playground", "secrets.kdbx")

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("playground database missing at %s: %v", path, err)
	}

	return path
}

func mkValue(key, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value},
	}
}

func mkProtectedValue(key, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key: key,
		Value: gokeepasslib.V{
			Content:   value,
			Protected: w.NewBoolWrapper(true),
		},
	}
}

func mkEntry(title, password string) gokeepasslib.Entry {
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, mkValue("Title", title))
	entry.Values = append(entry.Values, mkProtectedValue("Password", password))
	return entry
}

func createUnitKeePassDB(t *testing.T) string {
	t.Helper()

	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "Secrets"

	rootGroup.Entries = append(rootGroup.Entries, mkEntry("DB_PASSWORD", "db-pass"))
	rootGroup.Entries = append(rootGroup.Entries, mkEntry("API_KEY", "api-key"))
	rootGroup.Entries = append(rootGroup.Entries, mkEntry("NGINX_PASSWORD", "nginx-pass"))

	testGroup := gokeepasslib.NewGroup()
	testGroup.Name = "test"
	testGroup.Entries = append(testGroup.Entries, mkEntry("test1", "value-1"))
	testGroup.Entries = append(testGroup.Entries, mkEntry("test2", "value-2"))

	emptyNested := gokeepasslib.NewGroup()
	emptyNested.Name = "empty-nested"
	testGroup.Groups = append(testGroup.Groups, emptyNested)

	emptyGroup := gokeepasslib.NewGroup()
	emptyGroup.Name = "empty-group"

	rootGroup.Groups = append(rootGroup.Groups, testGroup)
	rootGroup.Groups = append(rootGroup.Groups, emptyGroup)

	db := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(testDBPassword),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}

	db.LockProtectedEntries()

	dbPath := filepath.Join(t.TempDir(), "unit-secrets.kdbx")

	file, err := os.Create(dbPath)
	if err != nil {
		t.Fatalf("create db: %v", err)
	}
	defer file.Close()

	if err := gokeepasslib.NewEncoder(file).Encode(db); err != nil {
		t.Fatalf("encode db: %v", err)
	}

	return dbPath
}

func TestLoadKeePass_Playground(t *testing.T) {
	dbPath := playgroundDBPath(t)

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
		"Secrets/API_KEY":        "api-key",
		"Secrets/DB_PASSWORD":    "db-pass",
		"Secrets/NGINX_PASSWORD": "nginx-pass",
		"Secrets/test/test1":     "value-1",
		"Secrets/test/test2":     "value-2",
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
		"DB_PASSWORD":    "db-pass",
		"API_KEY":        "api-key",
		"NGINX_PASSWORD": "nginx-pass",
		"test1":          "value-1",
		"test2":          "value-2",
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
