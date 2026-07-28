package sato

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// ReadPassword reads the KeePass master password from stdin.
// In interactive mode, it hides the input; in piped mode, it reads the password directly.
func ReadPassword() (string, error) {
	fmt.Fprint(os.Stderr, "KeePass password: ")

	// Check if stdin is a terminal (interactive mode)
	if term.IsTerminal(int(syscall.Stdin)) {
		pass, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", err
		}
		return string(pass), nil
	}

	// Piped input: read password directly from stdin
	reader := bufio.NewReader(os.Stdin)
	pass, err := reader.ReadString('\n')
	if err != nil && len(pass) == 0 {
		return "", err
	}

	return strings.TrimSpace(pass), nil
}
