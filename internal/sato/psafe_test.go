package sato

import (
	"path/filepath"
	"testing"

	"github.com/tkuhlman/gopwsafe/pwsafe"
)

func createPasswordSafeTestDB(t *testing.T, extension string) string {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "secrets"+extension)

	db := pwsafe.NewV3("SATO test database", "sato")

	db.SetRecord(pwsafe.Record{
		Title:    "DB_PASSWORD",
		Password: "super_secret_db_password_123",
		Group:    "Secrets.Database",
	})

	db.SetRecord(pwsafe.Record{
		Title:    "API_KEY",
		Password: "sk-test-api-key-xyz789",
		Group:    "Secrets.API",
	})

	if err := pwsafe.WritePWSafeFile(db, dbPath); err != nil {
		t.Fatalf("write Password Safe database: %v", err)
	}

	return dbPath
}

func TestLoadPasswordSafeData_PSafe3(t *testing.T) {
	dbPath := createPasswordSafeTestDB(t, ".psafe3")

	data, err := LoadPasswordSafeData(dbPath, "sato")
	if err != nil {
		t.Fatalf("LoadPasswordSafeData: %v", err)
	}

	values := make(map[string]string)
	paths := make(map[string]string)

	for _, entry := range data.Entries {
		values[entry.Name] = entry.Value
		paths[entry.Name] = entry.Path
	}

	if values["DB_PASSWORD"] != "super_secret_db_password_123" {
		t.Fatalf("unexpected DB_PASSWORD value: %q", values["DB_PASSWORD"])
	}

	if values["API_KEY"] != "sk-test-api-key-xyz789" {
		t.Fatalf("unexpected API_KEY value: %q", values["API_KEY"])
	}

	if paths["DB_PASSWORD"] != "Secrets/Database/DB_PASSWORD" {
		t.Fatalf("unexpected DB_PASSWORD path: %q", paths["DB_PASSWORD"])
	}

	if paths["API_KEY"] != "Secrets/API/API_KEY" {
		t.Fatalf("unexpected API_KEY path: %q", paths["API_KEY"])
	}
}

func TestLoadPasswordSafeData_IBak(t *testing.T) {
	dbPath := createPasswordSafeTestDB(t, ".ibak")

	data, err := LoadPasswordSafeData(dbPath, "sato")
	if err != nil {
		t.Fatalf("LoadPasswordSafeData: %v", err)
	}

	if len(data.Entries) != 2 {
		t.Fatalf("entries count = %d; want 2", len(data.Entries))
	}
}

func TestLoadPasswordSafeData_WrongPassword(t *testing.T) {
	dbPath := createPasswordSafeTestDB(t, ".psafe3")

	if _, err := LoadPasswordSafeData(dbPath, "wrong-password"); err == nil {
		t.Fatal("expected error for wrong Password Safe password")
	}
}

func TestLoadSecretStore_PSafe3(t *testing.T) {
	dbPath := createPasswordSafeTestDB(t, ".psafe3")

	data, err := LoadSecretStore(dbPath, "sato")
	if err != nil {
		t.Fatalf("LoadSecretStore: %v", err)
	}

	if len(data.Entries) != 2 {
		t.Fatalf("entries count = %d; want 2", len(data.Entries))
	}
}

func TestLoadSecretStore_IBak(t *testing.T) {
	dbPath := createPasswordSafeTestDB(t, ".ibak")

	data, err := LoadSecretStore(dbPath, "sato")
	if err != nil {
		t.Fatalf("LoadSecretStore: %v", err)
	}

	if len(data.Entries) != 2 {
		t.Fatalf("entries count = %d; want 2", len(data.Entries))
	}
}

func TestLoadSecretStore_UnsupportedFormat(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "secrets.unknown")

	_, err := LoadSecretStore(dbPath, "sato")
	if err == nil {
		t.Fatal("expected unsupported database format error")
	}
}
