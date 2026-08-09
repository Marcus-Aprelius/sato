package sato

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var allowedGitCommands = map[string]bool{
	"clone": true,
	"fetch": true,
	"pull":  true,
	"push":  true,
}

func IsAllowedGitCommand(args []string) bool {
	if len(args) < 2 {
		return false
	}

	if args[0] != "git" {
		return false
	}

	return allowedGitCommands[args[1]]
}

func repoSecretName(repoURL string) string {
	u, err := url.Parse(repoURL)
	if err == nil && u.Path != "" {
		base := filepath.Base(u.Path)
		return strings.TrimSuffix(base, ".git")
	}

	base := filepath.Base(repoURL)
	return strings.TrimSuffix(base, ".git")
}

func gitRemoteURL() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git remote origin URL: %w", err)
	}

	remoteURL := strings.TrimSpace(string(out))
	if remoteURL == "" {
		return "", fmt.Errorf("git remote origin URL is empty")
	}

	return remoteURL, nil
}

func gitSecretFromArgs(args []string, secretOverride string) (string, []string, error) {
	if len(args) < 2 {
		return "", nil, fmt.Errorf("usage: sato [--secret NAME] git clone|fetch|pull|push")
	}

	if secretOverride != "" {
		return secretOverride, args, nil
	}

	switch args[1] {
	case "clone":
		if len(args) < 3 {
			return "", nil, fmt.Errorf("usage: sato [--secret NAME] git clone URL")
		}

		return repoSecretName(args[2]), args, nil

	case "fetch", "pull", "push":
		remoteURL, err := gitRemoteURL()
		if err != nil {
			return "", nil, err
		}

		return repoSecretName(remoteURL), args, nil

	default:
		return "", nil, fmt.Errorf("unsupported git command: %s", args[1])
	}
}

func RunGitCommand(args []string, token string) error {
	tmpFile, err := os.CreateTemp("", "sato-git-askpass-*")
	if err != nil {
		return err
	}

	askpassPath := tmpFile.Name()

	askpassScript := `#!/usr/bin/env sh
case "$1" in
  *Username*) echo "oauth2" ;;
  *Password*) echo "$SATO_GIT_TOKEN" ;;
  *) echo "$SATO_GIT_TOKEN" ;;
esac
`

	if _, err := tmpFile.WriteString(askpassScript); err != nil {
		tmpFile.Close()
		os.Remove(askpassPath)
		return err
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(askpassPath)
		return err
	}

	if err := os.Chmod(askpassPath, 0o700); err != nil {
		os.Remove(askpassPath)
		return err
	}

	defer os.Remove(askpassPath)

	gitArgs := append(
		[]string{
			"-c", "credential.helper=",
		},
		args[1:]...,
	)

	cmd := exec.Command("git", gitArgs...)

	cmd.Env = append(
		buildCommandEnv(nil),
		"GIT_ASKPASS="+askpassPath,
		"GIT_TERMINAL_PROMPT=0",
		"SATO_GIT_TOKEN="+token,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
