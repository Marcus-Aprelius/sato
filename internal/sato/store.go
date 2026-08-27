package sato

import (
	"fmt"
	"path/filepath"
	"strings"
)

// LoadSecretStore loads secrets from a supported database format.
func LoadSecretStore(dbPath string, password string) (KeePassData, error) {
	extension := strings.ToLower(filepath.Ext(dbPath))

	switch extension {
	case ".kdbx":
		return LoadKeePassData(dbPath, password)

	case ".psafe3", ".ibak":
		return LoadPasswordSafeData(dbPath, password)

	default:
		return KeePassData{}, fmt.Errorf(
			"unsupported database format %q; supported: .kdbx, .psafe3, .ibak",
			extension,
		)
	}
}

// LoadSecretEnvironment loads secrets from any supported database
// and converts them to environment variables.
func LoadSecretEnvironment(dbPath string, password string) (map[string]string, error) {
	data, err := LoadSecretStore(dbPath, password)
	if err != nil {
		return nil, err
	}

	env := make(map[string]string, len(data.Entries))

	for _, entry := range data.Entries {
		if entry.Name != "" {
			env[entry.Name] = entry.Value
		}
	}

	return env, nil
}
