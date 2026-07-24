package main

import (
	"fmt"
	"os"

	"sato/internal/sato"
)

func main() {
	if err := sato.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}
