package sato

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/tkuhlman/gopwsafe/pwsafe"
)

// LoadPasswordSafeData opens a Password Safe V3 database.
// Both .psafe3 and .ibak files are handled as Password Safe V3 databases.
func LoadPasswordSafeData(dbPath string, masterPassword string) (KeePassData, error) {
	db, err := pwsafe.OpenPWSafeFile(dbPath, masterPassword)
	if err != nil {
		return KeePassData{}, fmt.Errorf(
			"failed to open Password Safe database: %w",
			err,
		)
	}

	data := KeePassData{
		Entries:         make([]SecretEntry, 0, len(db.Records)),
		EmptyGroupPaths: make([]string, 0),
	}

	for _, record := range db.Records {
		name := strings.TrimSpace(record.Title)
		if name == "" {
			continue
		}

		group := normalizePasswordSafeGroup(record.Group)
		entryPath := name

		if group != "" {
			entryPath = path.Join(group, name)
		}

		data.Entries = append(data.Entries, SecretEntry{
			Name:  name,
			Value: record.Password,
			Path:  entryPath,
		})
	}

	sort.Slice(data.Entries, func(i, j int) bool {
		return data.Entries[i].Path < data.Entries[j].Path
	})

	return data, nil
}

func normalizePasswordSafeGroup(group string) string {
	group = strings.TrimSpace(group)
	group = strings.ReplaceAll(group, `\`, "/")
	group = strings.ReplaceAll(group, ".", "/")

	parts := strings.Split(group, "/")
	cleanParts := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	return strings.Join(cleanParts, "/")
}
