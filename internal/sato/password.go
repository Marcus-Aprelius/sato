package sato

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// ReadPassword reads the KeePass master password from stdin.
// The password prompt is printed to stderr.
func ReadPassword() (string, error) {
	return ReadPasswordWithPrompt(true)
}

// ReadPasswordWithPrompt reads the KeePass master password from stdin.
// In interactive mode, input is hidden.
// In piped mode, all stdin is read, so both "echo sato" and "printf sato" work.
func ReadPasswordWithPrompt(showPrompt bool) (string, error) {
	if showPrompt {
		fmt.Fprint(os.Stderr, "KeePass password: ")
	}

	if term.IsTerminal(int(syscall.Stdin)) {
		pass, err := term.ReadPassword(int(syscall.Stdin))
		if showPrompt {
			fmt.Fprintln(os.Stderr)
		}

		if err != nil {
			return "", err
		}

		return string(pass), nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	password := strings.TrimSpace(string(data))
	if password == "" {
		return "", io.EOF
	}

	return password, nil
}
