package sato

import (
	"os"
	"os/exec"
)

// RunCommand executes a command with environment variables from KeePass.
// Only KeePass variables are passed to the subprocess (strict isolation).
// Variables are not exported to current shell.
func RunCommand(cmd string, args []string, envVars map[string]string) error {
	command := exec.Command(cmd, args...)

	// Use only KeePass variables (strict isolation)
	command.Env = make([]string, 0, len(envVars))
	for k, v := range envVars {
		command.Env = append(command.Env, k+"="+v)
	}

	// Attach to current terminal for interactive I/O
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	return command.Run()
}
