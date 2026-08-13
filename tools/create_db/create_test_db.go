package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"sato/internal/testdb"

	"github.com/tkuhlman/gopwsafe/pwsafe"
)

const databasePassword = "sato"

func main() {
	if err := testdb.WriteDemoDatabase("secrets.kdbx", databasePassword); err != nil {
		log.Fatalf("create secrets.kdbx: %v", err)
	}

	if err := writePasswordSafeDatabase("secrets.psafe3"); err != nil {
		log.Fatalf("create secrets.psafe3: %v", err)
	}

	if err := writePasswordSafeDatabase("secrets.ibak"); err != nil {
		log.Fatalf("create secrets.ibak: %v", err)
	}

	fmt.Println("✅ Test databases created:")
	fmt.Println("   - secrets.kdbx")
	fmt.Println("   - secrets.psafe3")
	fmt.Println("   - secrets.ibak")
	fmt.Println()
	fmt.Println("   Password: sato")
	fmt.Println("   Entries:")
	fmt.Println("   - DB_PASSWORD")
	fmt.Println("   - API_KEY")
	fmt.Println("   - NGINX_PASSWORD")
	fmt.Println("   - test/test1")
	fmt.Println("   - test/test2")
	fmt.Println("   Empty groups in KDBX:")
	fmt.Println("   - empty-group")
	fmt.Println("   - test/empty-nested")
}

func writePasswordSafeDatabase(filePath string) error {
	db := pwsafe.NewV3("SATO test database", databasePassword)

	records := []pwsafe.Record{
		{
			Title:    "DB_PASSWORD",
			Password: "super_secret_db_password_123",
			Group:    "Secrets",
		},
		{
			Title:    "API_KEY",
			Password: "sk-test-api-key-xyz789",
			Group:    "Secrets",
		},
		{
			Title:    "NGINX_PASSWORD",
			Password: "nginx_password_demo_456",
			Group:    "Secrets",
		},
		{
			Title:    "test1",
			Password: "test_secret_value_1",
			Group:    "Secrets.test",
		},
		{
			Title:    "test2",
			Password: "test_secret_value_2",
			Group:    "Secrets.test",
		},
	}

	for _, record := range records {
		if _, err := rand.Read(record.UUID[:]); err != nil {
			return fmt.Errorf("generate record UUID: %w", err)
		}

		db.SetRecord(record)
	}

	if err := pwsafe.WritePWSafeFile(db, filePath); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
	}

	return nil
}
