package sato_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sato/internal/sato"
)

func captureCompletion(t *testing.T, args ...string) string {
	t.Helper()

	oldStdout := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stdout = writer

	sato.RunCompletion("bash", args...)

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	os.Stdout = oldStdout

	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}

	return string(out)
}

func TestCompletionContainsTopLevelCommands(t *testing.T) {
	out := captureCompletion(t)

	want := []string{
		"docker version completion help get git",
		"--db-path --secret",
		"complete -F _sato_completion sato",
	}

	for _, item := range want {
		if !strings.Contains(out, item) {
			t.Fatalf("completion missing %q", item)
		}
	}
}

func TestCompletionContainsGetCommandsAndFlags(t *testing.T) {
	out := captureCompletion(t)

	want := []string{
		"compgen -W \"secret secrets\"",
		"compgen -W \"--tree --show-empty-groups\"",
		"compgen -W \"-q --quiet\"",
	}

	for _, item := range want {
		if !strings.Contains(out, item) {
			t.Fatalf("completion missing %q", item)
		}
	}
}

func TestCompletionContainsGitCommands(t *testing.T) {
	out := captureCompletion(t)

	want := []string{
		"compgen -W \"clone fetch pull push help\"",
	}

	for _, item := range want {
		if !strings.Contains(out, item) {
			t.Fatalf("completion missing git item %q", item)
		}
	}
}

func TestCompletionContainsCompletionActions(t *testing.T) {
	out := captureCompletion(t)

	want := []string{
		"compgen -W \"bash\"",
		"compgen -W \"add delete update status\"",
	}

	for _, item := range want {
		if !strings.Contains(out, item) {
			t.Fatalf("completion missing completion action %q", item)
		}
	}
}

func TestCompletionContainsDockerDelegation(t *testing.T) {
	out := captureCompletion(t)

	want := []string{
		"_sato_docker_completion",
		"declare -F _docker",
		"COMP_WORDS=(\"docker\" \"${old_words[@]:2}\")",
		"_docker",
	}

	for _, item := range want {
		if !strings.Contains(out, item) {
			t.Fatalf("completion missing docker delegation %q", item)
		}
	}
}

func TestCompletionDoesNotContainKnownCorruptionMarkers(t *testing.T) {
	out := captureCompletion(t)

	bad := []string{
		"elif*",
		"*OMPREPLY",
		"*OM*REPLY",
		"--show-empty-groups*",
		"\"*cur\"",
	}

	for _, item := range bad {
		if strings.Contains(out, item) {
			t.Fatalf("completion contains corruption marker %q", item)
		}
	}
}

func TestCompletionBashStatusDisabledByDefault(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	out := captureCompletion(t, "status")

	if strings.TrimSpace(out) != "SATO bash completion: Disabled" {
		t.Fatalf("status = %q; want SATO bash completion: Disabled", strings.TrimSpace(out))
	}
}

func TestCompletionBashAddStatusDelete(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	addOut := captureCompletion(t, "add")
	if !strings.Contains(addOut, "SATO bash completion: Added") {
		t.Fatalf("add output unexpected: %q", addOut)
	}

	if !strings.Contains(addOut, "Reload shell:") {
		t.Fatalf("add output missing reload instruction: %q", addOut)
	}

	if !strings.Contains(addOut, "source ~/.bashrc") {
		t.Fatalf("add output missing source command: %q", addOut)
	}

	bashrcPath := filepath.Join(home, ".bashrc")
	completionPath := filepath.Join(home, ".sato-completion.bash")

	bashrcData, err := os.ReadFile(bashrcPath)
	if err != nil {
		t.Fatalf("read .bashrc: %v", err)
	}

	if !strings.Contains(string(bashrcData), "# sato completion start") {
		t.Fatalf(".bashrc missing start marker:\n%s", string(bashrcData))
	}

	if !strings.Contains(string(bashrcData), ".sato-completion.bash") {
		t.Fatalf(".bashrc missing completion source:\n%s", string(bashrcData))
	}

	if _, err := os.Stat(completionPath); err != nil {
		t.Fatalf("completion file missing: %v", err)
	}

	statusOut := captureCompletion(t, "status")
	if strings.TrimSpace(statusOut) != "SATO bash completion: Enabled" {
		t.Fatalf("status = %q; want SATO bash completion: Enabled", strings.TrimSpace(statusOut))
	}

	deleteOut := captureCompletion(t, "delete")
	if !strings.Contains(deleteOut, "SATO bash completion: Removed") {
		t.Fatalf("delete output unexpected: %q", deleteOut)
	}

	if !strings.Contains(deleteOut, "Reload shell:") {
		t.Fatalf("delete output missing reload instruction: %q", deleteOut)
	}

	if !strings.Contains(deleteOut, "source ~/.bashrc") {
		t.Fatalf("delete output missing source command: %q", deleteOut)
	}

	statusOut = captureCompletion(t, "status")
	if strings.TrimSpace(statusOut) != "SATO bash completion: Disabled" {
		t.Fatalf("status after delete = %q; want SATO bash completion: Disabled", strings.TrimSpace(statusOut))
	}
}

func TestCompletionBashUpdateCreatesCompletionFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	out := captureCompletion(t, "update")
	if !strings.Contains(out, "SATO bash completion: Updated") {
		t.Fatalf("update output unexpected: %q", out)
	}

	if !strings.Contains(out, "Reload shell:") {
		t.Fatalf("update output missing reload instruction: %q", out)
	}

	if !strings.Contains(out, "source ~/.bashrc") {
		t.Fatalf("update output missing source command: %q", out)
	}

	completionPath := filepath.Join(home, ".sato-completion.bash")

	data, err := os.ReadFile(completionPath)
	if err != nil {
		t.Fatalf("read completion file: %v", err)
	}

	if !strings.Contains(string(data), "complete -F _sato_completion sato") {
		t.Fatalf("completion file missing complete command")
	}
}
