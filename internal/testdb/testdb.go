package testdb

import (
	"os"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

const DefaultPassword = "sato"

func Value(key, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value},
	}
}

func ProtectedValue(key, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key: key,
		Value: gokeepasslib.V{
			Content:   value,
			Protected: w.NewBoolWrapper(true),
		},
	}
}

func Entry(title, password string) gokeepasslib.Entry {
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, Value("Title", title))
	entry.Values = append(entry.Values, ProtectedValue("Password", password))
	return entry
}

func NewDemoDatabase(password string) *gokeepasslib.Database {
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "Secrets"

	rootGroup.Entries = append(rootGroup.Entries, Entry("DB_PASSWORD", "super_secret_db_password_123"))
	rootGroup.Entries = append(rootGroup.Entries, Entry("API_KEY", "sk-test-api-key-xyz789"))
	rootGroup.Entries = append(rootGroup.Entries, Entry("NGINX_PASSWORD", "nginx_password_demo_456"))

	testGroup := gokeepasslib.NewGroup()
	testGroup.Name = "test"
	testGroup.Entries = append(testGroup.Entries, Entry("test1", "test_secret_value_1"))
	testGroup.Entries = append(testGroup.Entries, Entry("test2", "test_secret_value_2"))

	emptyNestedGroup := gokeepasslib.NewGroup()
	emptyNestedGroup.Name = "empty-nested"
	testGroup.Groups = append(testGroup.Groups, emptyNestedGroup)

	emptyGroup := gokeepasslib.NewGroup()
	emptyGroup.Name = "empty-group"

	rootGroup.Groups = append(rootGroup.Groups, testGroup)
	rootGroup.Groups = append(rootGroup.Groups, emptyGroup)

	db := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(password),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}

	db.LockProtectedEntries()

	return db
}

func WriteDemoDatabase(path string, password string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return gokeepasslib.NewEncoder(file).Encode(NewDemoDatabase(password))
}
