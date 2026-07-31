package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"sato/internal/sato"
)

func main() {
	if err := sato.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			code := exitErr.ExitCode()
			if code != 0 {
				os.Exit(code)
			}
		}

		os.Exit(1)
	}
}
