package sato

import (
	"errors"
	"fmt"
	"strings"
)

func emptyGroupDisplayPaths(entries []SecretEntry, emptyGroupPaths []string) []string {
	allPaths := make([][]string, 0, len(entries)+len(emptyGroupPaths))
	groupPaths := make([][]string, 0, len(emptyGroupPaths))

	for _, entry := range entries {
		if parts := cleanPathParts(entry.Path); len(parts) > 0 {
			allPaths = append(allPaths, parts)
		}
	}

	for _, groupPath := range emptyGroupPaths {
		if parts := cleanPathParts(groupPath); len(parts) > 0 {
			groupPaths = append(groupPaths, parts)
			allPaths = append(allPaths, parts)
		}
	}

	if shouldStripTopGroup(allPaths) {
		for i := range groupPaths {
			groupPaths[i] = groupPaths[i][1:]
		}
	}

	values := make([]string, 0, len(groupPaths))

	for _, parts := range groupPaths {
		if len(parts) == 0 {
			continue
		}

		values = append(values, strings.Join(parts, "/")+"/")
	}

	return values
}

func parseGetSecretArgs(args []string, inheritedQuiet bool) (string, bool, error) {
	secretName := ""
	quiet := inheritedQuiet

	for _, arg := range args {
		switch arg {
		case "-q", "--quiet":
			quiet = true
		default:
			if strings.HasPrefix(arg, "-") {
				return "", quiet, fmt.Errorf("unknown option: %s", arg)
			}

			if secretName != "" {
				return "", quiet, errors.New("usage: sato get secret NAME [-q|--quiet]")
			}

			secretName = arg
		}
	}

	if secretName == "" {
		return "", quiet, errors.New("usage: sato get secret NAME [-q|--quiet]")
	}

	return secretName, quiet, nil
}

func findSecretValue(entries []SecretEntry, name string) (string, error) {
	matches := make([]SecretEntry, 0)

	for _, entry := range entries {
		if entry.Name == name {
			matches = append(matches, entry)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("secret not found: %s", name)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("secret name is not unique: %s", name)
	}

	return matches[0].Value, nil
}
