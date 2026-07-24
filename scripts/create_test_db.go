package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

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

func main() {
	// Create root group
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "Secrets"

	// Root-level entries
	rootGroup.Entries = append(rootGroup.Entries, mkEntry("DB_PASSWORD", "super_secret_db_password_123"))
	rootGroup.Entries = append(rootGroup.Entries, mkEntry("API_KEY", "sk-test-api-key-xyz789"))
	rootGroup.Entries = append(rootGroup.Entries, mkEntry("NGINX_PASSWORD", "nginx_password_demo_456"))

	// Nested group with entries
	testGroup := gokeepasslib.NewGroup()
	testGroup.Name = "test"
	testGroup.Entries = append(testGroup.Entries, mkEntry("test1", "test_secret_value_1"))
	testGroup.Entries = append(testGroup.Entries, mkEntry("test2", "test_secret_value_2"))

	// Empty nested group
	emptyNestedGroup := gokeepasslib.NewGroup()
	emptyNestedGroup.Name = "empty-nested"
	testGroup.Groups = append(testGroup.Groups, emptyNestedGroup)

	// Empty root-level group
	emptyGroup := gokeepasslib.NewGroup()
	emptyGroup.Name = "empty-group"

	// Attach groups
	rootGroup.Groups = append(rootGroup.Groups, testGroup)
	rootGroup.Groups = append(rootGroup.Groups, emptyGroup)

	// Create database
	db := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials("sato"),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}

	// Lock protected entries
	db.LockProtectedEntries()

	// Create and encode to file
	file, err := os.Create("secrets.kdbx")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := gokeepasslib.NewEncoder(file).Encode(db); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Test database created: secrets.kdbx")
	fmt.Println("   Password: sato")
	fmt.Println("   Entries:")
	fmt.Println("   - DB_PASSWORD")
	fmt.Println("   - API_KEY")
	fmt.Println("   - NGINX_PASSWORD")
	fmt.Println("   - test/test1")
	fmt.Println("   - test/test2")
	fmt.Println("   Empty groups:")
	fmt.Println("   - empty-group")
	fmt.Println("   - test/empty-nested")
}
