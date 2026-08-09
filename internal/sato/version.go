package sato

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Version and Build info.
// These can be set during build via ldflags.
var (
	Version   = "v0.0.0-dev"
	GitCommit = "unknown"
)

func RunVersion() {
	fmt.Printf("%s\n", Version)
	fmt.Printf("[SHA: %s]\n", GitCommit)
	fmt.Println()

	if version := gitVersion(); version != "" {
		fmt.Printf("Git: %s\n", version)
	}
	if version := dockerVersion(); version != "" {
		fmt.Printf("Docker: %s\n", version)
	}
	if version := dockerComposeVersion(); version != "" {
		fmt.Printf("Docker Compose: %s\n", version)
	}
	if version := osVersion(); version != "" {
		fmt.Printf("OS: %s\n", version)
	}
}

func osVersion() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[len("PRETTY_NAME="):], `"`)
		}
	}

	return ""
}

func dockerVersion() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "docker", "version", "--format", "{{.Client.Version}}").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func dockerComposeVersion() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(
		ctx,
		"docker", "compose", "version", "--short",
	).Output()
	if err != nil {
		return ""
	}

	version := strings.TrimSpace(string(out))
	if version == "" {
		return ""
	}

	if i := strings.Index(version, "+"); i >= 0 {
		version = version[:i]
	}

	return version
}

func gitVersion() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "git", "--version").Output()
	if err != nil {
		return ""
	}

	fields := strings.Fields(string(out))
	if len(fields) < 3 {
		return ""
	}

	return fields[2]
}
