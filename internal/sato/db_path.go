package sato

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// FindKDBX returns the first *.kdbx file found in dir, or "".
func FindKDBX(dir string) string {
	matches, err := filepath.Glob(filepath.Join(dir, "*.kdbx"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	return matches[0]
}

// FindDBPath resolves the KeePass-compatible .kdbx database path using the priority order:
//  1. --db-path flag
//  2. ~/.sato/*.kdbx
//  3. SATO_DB_PATH environment variable
//
// Returns "" if none is available.
func FindDBPath(dbPathFlag string) string {
	if dbPathFlag != "" {
		if _, err := os.Stat(dbPathFlag); err == nil {
			return dbPathFlag
		}
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		if p := FindKDBX(filepath.Join(homeDir, ".sato")); p != "" {
			return p
		}
	}

	if envPath := os.Getenv("SATO_DB_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	return ""
}

func fileAccess(path string) string {
	readOK := unix.Access(path, unix.R_OK) == nil
	writeOK := unix.Access(path, unix.W_OK) == nil

	switch {
	case readOK && writeOK:
		return "RW"
	case readOK:
		return "RO"
	default:
		return "--"
	}
}

func satoLocalStatus() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	if p := FindKDBX(filepath.Join(homeDir, ".sato")); p != "" {
		return p
	}
	return ""
}

func printStatus(dbPathFlag string) {
	printTitle()
	fmt.Println()

	activePriority := 0

	dbPath := FindDBPath(dbPathFlag)

	localPath := satoLocalStatus()
	envPath := os.Getenv("SATO_DB_PATH")

	if dbPath != "" {
		if dbPathFlag != "" && dbPath == dbPathFlag {
			activePriority = 1
		} else if localPath != "" && dbPath == localPath {
			activePriority = 2
		} else if envPath != "" && dbPath == envPath {
			activePriority = 3
		}
	}

	dbPathStatus := "not set"

	if dbPathFlag != "" {
		if _, err := os.Stat(dbPathFlag); err == nil {
			dbPathStatus = "set"
		} else {
			dbPathStatus = "set (DB is missing)"
		}
	}

	localStatus := "not found"

	if localPath != "" {
		localStatus = "found"
	}

	envStatus := "not set"

	if envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			envStatus = "set"
		} else {
			envStatus = "set (DB is missing)"
		}
	}

	priCell := func(n int) string {
		if n == activePriority {
			return fmt.Sprintf("[%d]", n)
		}
		return fmt.Sprintf(" %d", n)
	}

	modeCell := func(path string) string {
		if path == "" {
			return ""
		}

		if _, err := os.Stat(path); err != nil {
			return ""
		}

		mode := fileAccess(path)
		if mode == "--" {
			return ""
		}

		return mode
	}

	showMode := dbPath != ""

	if showMode {
		fmt.Println(" Priority | Path           | Status              | Mode")
		fmt.Println("----------|----------------|---------------------|------")

		fmt.Printf("%-9s | %-14s | %-19s | %s\n",
			priCell(1),
			"--db-path",
			dbPathStatus,
			modeCell(dbPathFlag),
		)

		fmt.Printf("%-9s | %-14s | %-19s | %s\n",
			priCell(2),
			"~/.sato/*.kdbx",
			localStatus,
			modeCell(localPath),
		)

		fmt.Printf("%-9s | %-14s | %-19s | %s\n",
			priCell(3),
			"SATO_DB_PATH",
			envStatus,
			modeCell(envPath),
		)
	} else {
		fmt.Println(" Priority | Path           | Status")
		fmt.Println("----------|----------------|-----------")

		fmt.Printf("%-9s | %-14s | %s\n",
			priCell(1),
			"--db-path",
			dbPathStatus,
		)

		fmt.Printf("%-9s | %-14s | %s\n",
			priCell(2),
			"~/.sato/*.kdbx",
			localStatus,
		)

		fmt.Printf("%-9s | %-14s | %s\n",
			priCell(3),
			"SATO_DB_PATH",
			envStatus,
		)
	}

	fmt.Println()

	if dbPath != "" {
		fmt.Printf("Database: %s\n", dbPath)
	} else {
		fmt.Println("Database: not found")
	}
}
