package sato_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"sato/internal/sato"
)

func captureCompletion(t *testing.T) string {
	t.Helper()

	oldStdout := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stdout = writer

	sato.RunCompletion("bash")

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
		"docker version completion help get",
		"--db-path",
		"complete -F _sato_completion sato",
	}

	for _, item := range want {
		if !strings.Contains(out, item) {
			t.Fatalf("completion missing %q", item)
		}
	}
}

func TestCompletionContainsGetSecretsFlags(t *testing.T) {
	out := captureCompletion(t)

	want := []string{
		"compgen -W \"secrets\"",
		"compgen -W \"--tree --show-empty-groups\"",
	}

	for _, item := range want {
		if !strings.Contains(out, item) {
			t.Fatalf("completion missing %q", item)
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
