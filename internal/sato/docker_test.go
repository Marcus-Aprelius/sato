package sato

import (
	"strings"
	"testing"
)

func TestBuildCommandEnv_AllowsOnlySafeParentEnvAndKeePassSecrets(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")
	t.Setenv("HOME", "/home/test")
	t.Setenv("DOCKER_HOST", "ssh://user@example")
	t.Setenv("LEAK_ME", "must-not-leak")

	env := buildCommandEnv(map[string]string{
		"DB_PASSWORD": "secret-db-password",
		"API_KEY":     "secret-api-key",
	})

	got := strings.Join(env, "\n")

	wantContains := []string{
		"PATH=/usr/bin",
		"HOME=/home/test",
		"DOCKER_HOST=ssh://user@example",
		"DB_PASSWORD=secret-db-password",
		"API_KEY=secret-api-key",
	}

	for _, item := range wantContains {
		if !strings.Contains(got, item) {
			t.Fatalf("env missing %q; env:\n%s", item, got)
		}
	}

	if strings.Contains(got, "LEAK_ME=must-not-leak") {
		t.Fatalf("unexpected parent env leak; env:\n%s", got)
	}
}
