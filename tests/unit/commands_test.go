package sato_test

import (
	"testing"

	"sato/internal/sato"
)

// TestIsAllowedCommand is security-critical:
// only "docker compose" commands may be executed by SATO through docker path.
func TestIsAllowedCommand(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		allowed bool
	}{
		{"docker compose up", []string{"docker", "compose", "up"}, true},
		{"docker compose config", []string{"docker", "compose", "config"}, true},
		{"docker compose watch", []string{"docker", "compose", "watch"}, true},
		{"docker compose custom command", []string{"docker", "compose", "future-command"}, true},
		{"docker compose with flags", []string{"docker", "compose", "up", "-d"}, true},
		{"docker compose no subcmd", []string{"docker", "compose"}, true},

		{"docker alone", []string{"docker"}, false},
		{"docker run", []string{"docker", "run"}, false},
		{"docker ps", []string{"docker", "ps"}, false},
		{"docker-compose legacy", []string{"docker-compose", "up"}, false},
		{"git clone is checked separately", []string{"git", "clone", "https://github.com/org/repo.git"}, false},
		{"kubectl", []string{"kubectl", "get", "pods"}, false},
		{"helm", []string{"helm", "list"}, false},
		{"shell", []string{"sh", "-c", "echo $DB_PASSWORD"}, false},
		{"bash", []string{"bash", "-c", "echo hi"}, false},
		{"deploy script", []string{"./deploy.sh"}, false},
		{"empty args", []string{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := sato.IsAllowedCommand(c.args)
			if got != c.allowed {
				t.Fatalf("IsAllowedCommand(%v) = %v; want %v", c.args, got, c.allowed)
			}
		})
	}
}

func TestHasComposeSubcommand(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{"docker compose up", []string{"docker", "compose", "up"}, true},
		{"docker compose config", []string{"docker", "compose", "config"}, true},
		{"docker compose watch", []string{"docker", "compose", "watch"}, true},
		{"docker compose up with flags", []string{"docker", "compose", "up", "-d"}, true},

		{"docker compose missing", []string{"docker", "compose"}, false},
		{"docker only", []string{"docker"}, false},
		{"empty", []string{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := sato.HasComposeSubcommand(c.args)
			if got != c.want {
				t.Fatalf("HasComposeSubcommand(%v) = %v; want %v", c.args, got, c.want)
			}
		})
	}
}

func TestIsAllowedGitCommand(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		allowed bool
	}{
		{"git clone", []string{"git", "clone", "https://github.com/org/repo.git"}, true},
		{"git fetch", []string{"git", "fetch"}, true},
		{"git pull", []string{"git", "pull"}, true},
		{"git push", []string{"git", "push"}, true},

		{"git alone", []string{"git"}, false},
		{"git help", []string{"git", "help"}, false},
		{"git status", []string{"git", "status"}, false},
		{"git config", []string{"git", "config"}, false},
		{"git commit", []string{"git", "commit"}, false},
		{"git checkout", []string{"git", "checkout"}, false},
		{"git remote", []string{"git", "remote", "-v"}, false},
		{"not git", []string{"kubectl", "get", "pods"}, false},
		{"empty args", []string{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := sato.IsAllowedGitCommand(c.args)
			if got != c.allowed {
				t.Fatalf("IsAllowedGitCommand(%v) = %v; want %v", c.args, got, c.allowed)
			}
		})
	}
}
