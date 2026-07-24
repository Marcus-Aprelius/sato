package sato

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
)

type SecretEntry struct {
	Name  string
	Value string
	Path  string
}

type KeePassData struct {
	Entries         []SecretEntry
	EmptyGroupPaths []string
}

// LoadKeePass opens a KeePass-compatible .kdbx database and extracts
// environment variables from all groups recursively.
// Returns a map of Title -> Password.
func LoadKeePass(dbPath string, masterPassword string) (map[string]string, error) {
	entries, err := LoadKeePassEntries(dbPath, masterPassword)
	if err != nil {
		return nil, err
	}

	env := make(map[string]string)

	for _, entry := range entries {
		if entry.Name != "" {
			env[entry.Name] = entry.Value
		}
	}

	return env, nil
}

// LoadKeePassEntries opens a KeePass-compatible .kdbx database and returns
// all secret entries from all groups recursively.
func LoadKeePassEntries(dbPath string, masterPassword string) ([]SecretEntry, error) {
	data, err := LoadKeePassData(dbPath, masterPassword)
	if err != nil {
		return nil, err
	}

	return data.Entries, nil
}

// LoadKeePassData opens a KeePass-compatible .kdbx database and returns
// all secret entries plus empty group paths.
func LoadKeePassData(dbPath string, masterPassword string) (KeePassData, error) {
	file, err := os.Open(dbPath)
	if err != nil {
		return KeePassData{}, fmt.Errorf("failed to open database: %w", err)
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(masterPassword)

	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		return KeePassData{}, fmt.Errorf("failed to decode database: %w", err)
	}

	if err := db.UnlockProtectedEntries(); err != nil {
		return KeePassData{}, fmt.Errorf("failed to unlock protected entries: %w", err)
	}

	data := KeePassData{
		Entries:         make([]SecretEntry, 0),
		EmptyGroupPaths: make([]string, 0),
	}

	if db.Content.Root == nil {
		return data, nil
	}

	for _, group := range db.Content.Root.Groups {
		walkKeePassGroup(group, "", &data)
	}

	return data, nil
}

func walkKeePassGroup(group gokeepasslib.Group, parentPath string, data *KeePassData) bool {
	groupPath := group.Name

	if parentPath != "" {
		groupPath = path.Join(parentPath, group.Name)
	}

	hasSecrets := false

	for _, entry := range group.Entries {
		title := ""
		password := ""

		for _, v := range entry.Values {
			switch v.Key {
			case "Title":
				title = strings.TrimSpace(v.Value.Content)
			case "Password":
				password = v.Value.Content
			}
		}

		if title != "" {
			hasSecrets = true

			data.Entries = append(data.Entries, SecretEntry{
				Name:  title,
				Value: password,
				Path:  path.Join(groupPath, title),
			})
		}
	}

	childHasSecrets := false

	for _, child := range group.Groups {
		if walkKeePassGroup(child, groupPath, data) {
			childHasSecrets = true
		}
	}

	if !hasSecrets && !childHasSecrets {
		data.EmptyGroupPaths = append(data.EmptyGroupPaths, groupPath)
	}

	return hasSecrets || childHasSecrets
}
