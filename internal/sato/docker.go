package sato

import (
	"os"
	"os/exec"
)

var allowedParentEnv = []string{
	"PATH",
	"HOME",
	"TERM",
	"DOCKER_HOST",
	"DOCKER_CONTEXT",
	"DOCKER_CONFIG",
	"DOCKER_TLS_VERIFY",
	"DOCKER_CERT_PATH",
	"SSH_AUTH_SOCK",
	"XDG_CONFIG_HOME",
}

func buildCommandEnv(envVars map[string]string) []string {
	env := make([]string, 0, len(allowedParentEnv)+len(envVars))

	for _, key := range allowedParentEnv {
		if value, ok := os.LookupEnv(key); ok {
			env = append(env, key+"="+value)
		}
	}

	for key, value := range envVars {
		env = append(env, key+"="+value)
	}

	return env
}

// RunCommand executes a command with environment variables from KeePass.
// Only selected safe parent variables and KeePass variables are passed.
// Variables are not exported to current shell.
func RunCommand(cmd string, args []string, envVars map[string]string) error {
	command := exec.Command(cmd, args...)

	command.Env = buildCommandEnv(envVars)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	return command.Run()
}
