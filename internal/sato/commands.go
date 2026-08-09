package sato

// IsAllowedCommand reports whether the args represent an allowed docker command.
// Only "docker compose" is permitted here.
func IsAllowedCommand(args []string) bool {
	if len(args) < 2 {
		return false
	}

	return args[0] == "docker" && args[1] == "compose"
}

func HasComposeSubcommand(args []string) bool {
	return len(args) >= 3
}
