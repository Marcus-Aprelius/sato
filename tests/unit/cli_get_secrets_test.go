package sato_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func runSATO(t *testing.T, args ...string) (string, string) {
	t.Helper()

	cmdArgs := append([]string{"run", "."}, args...)

	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = projectRoot(t)
	cmd.Stdin = strings.NewReader(testDBPassword + "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf(
			"go run . %v failed: %v\nstdout:\n%s\nstderr:\n%s",
			args,
			err,
			stdout.String(),
			stderr.String(),
		)
	}

	return stdout.String(), stderr.String()
}

func TestCLIGetSecrets_StdoutContainsOnlyNames(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	stdout, stderr := runSATO(t, "--db-path="+dbPath, "get", "secrets")

	if !strings.Contains(stderr, "[DB: "+dbPath+"]") {
		t.Fatalf("stderr does not contain db marker; stderr:\n%s", stderr)
	}

	if strings.Contains(stdout, "[DB:") {
		t.Fatalf("stdout must not contain db marker; stdout:\n%s", stdout)
	}

	want := []string{
		"API_KEY",
		"DB_PASSWORD",
		"NGINX_PASSWORD",
		"test1",
		"test2",
	}

	for _, item := range want {
		if !containsLine(stdout, item) {
			t.Fatalf("stdout missing %q; stdout:\n%s", item, stdout)
		}
	}
}

func TestCLIGetSecretsTree_HidesEmptyGroupsByDefault(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	stdout, _ := runSATO(t, "--db-path="+dbPath, "get", "secrets", "--tree")

	wantContains := []string{
		"├── API_KEY",
		"├── DB_PASSWORD",
		"├── NGINX_PASSWORD",
		"└── test/",
		"    ├── test1",
		"    └── test2",
	}

	for _, item := range wantContains {
		if !strings.Contains(stdout, item) {
			t.Fatalf("tree output missing %q; stdout:\n%s", item, stdout)
		}
	}

	mustNotContain := []string{
		"root",
		"Secrets/",
		"empty-group/",
		"empty-nested/",
	}

	for _, item := range mustNotContain {
		if strings.Contains(stdout, item) {
			t.Fatalf("tree output must not contain %q; stdout:\n%s", item, stdout)
		}
	}
}

func TestCLIGetSecretsTree_ShowEmptyGroups(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	stdout, _ := runSATO(
		t,
		"--db-path="+dbPath,
		"get",
		"secrets",
		"--tree",
		"--show-empty-groups",
	)

	wantContains := []string{
		"├── API_KEY",
		"├── DB_PASSWORD",
		"├── NGINX_PASSWORD",
		"├── empty-group/",
		"└── test/",
		"    ├── empty-nested/",
		"    ├── test1",
		"    └── test2",
	}

	for _, item := range wantContains {
		if !strings.Contains(stdout, item) {
			t.Fatalf("tree output missing %q; stdout:\n%s", item, stdout)
		}
	}
}

func TestCLIGetSecrets_ShowEmptyGroupsInListMode(t *testing.T) {
	dbPath := createUnitKeePassDB(t)

	stdout, _ := runSATO(
		t,
		"--db-path="+dbPath,
		"get",
		"secrets",
		"--show-empty-groups",
	)

	wantContains := []string{
		"API_KEY",
		"DB_PASSWORD",
		"NGINX_PASSWORD",
		"empty-group/",
		"test/empty-nested/",
	}

	for _, item := range wantContains {
		if !containsLine(stdout, item) {
			t.Fatalf("stdout missing %q; stdout:\n%s", item, stdout)
		}
	}

	if strings.Contains(stdout, "Secrets/") {
		t.Fatalf("common top-level group should be stripped; stdout:\n%s", stdout)
	}
}

func containsLine(output string, want string) bool {
	for _, line := range strings.Split(output, "\n") {
		if line == want {
			return true
		}
	}

	return false
}
