package sato

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// Version and Build info.
// These can be set during build via ldflags.
var (
	Version   = "v0.0.0-dev"
	GitCommit = "unknown"
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

func emptyGroupDisplayPaths(entries []SecretEntry, emptyGroupPaths []string) []string {
	allPaths := make([][]string, 0, len(entries)+len(emptyGroupPaths))
	groupPaths := make([][]string, 0, len(emptyGroupPaths))

	for _, entry := range entries {
		if parts := cleanPathParts(entry.Path); len(parts) > 0 {
			allPaths = append(allPaths, parts)
		}
	}

	for _, groupPath := range emptyGroupPaths {
		if parts := cleanPathParts(groupPath); len(parts) > 0 {
			groupPaths = append(groupPaths, parts)
			allPaths = append(allPaths, parts)
		}
	}

	if shouldStripTopGroup(allPaths) {
		for i := range groupPaths {
			groupPaths[i] = groupPaths[i][1:]
		}
	}

	values := make([]string, 0, len(groupPaths))

	for _, parts := range groupPaths {
		if len(parts) == 0 {
			continue
		}

		values = append(values, strings.Join(parts, "/")+"/")
	}

	return values
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

// IsAllowedCommand reports whether the args represent an allowed command.
// Only "docker compose" is permitted.
func IsAllowedCommand(args []string) bool {
	if len(args) < 2 {
		return false
	}
	return args[0] == "docker" && args[1] == "compose"
}

func HasComposeSubcommand(args []string) bool {
	return len(args) >= 3
}

func red(s string) string { return "\033[31m" + s + "\033[0m" }

func printTitle() {
	fmt.Printf(
		"%secure %sccess %sask %sperator - like sudo, but for secrets\n",
		red("S"),
		red("A"),
		red("T"),
		red("O"),
	)
}

func printLogoTitle() {
	fmt.Printf(`   _____      _______ ____  
  / ____|  /\|__   __| __ \    %secure
 | (___   /  \  | | | |  | |    %sccess
  \___ \ / /\ \ | | | |  | |     %sask
  ____) | ____ \| | | |__| |      %sperator
 |_____/_/    \_\_|  \____/    > like sudo, but for secrets

`,
		red("S"), red("A"), red("T"), red("O"),
	)
}

type treeNode struct {
	children map[string]*treeNode
	isGroup  bool
}

func newTreeNode() *treeNode {
	return &treeNode{
		children: make(map[string]*treeNode),
	}
}

func printSecretsTree(entries []SecretEntry, emptyGroupPaths []string, showEmptyGroups bool) {
	root := newTreeNode()

	secretPaths := make([][]string, 0, len(entries))
	groupPaths := make([][]string, 0, len(emptyGroupPaths))

	for _, entry := range entries {
		if parts := cleanPathParts(entry.Path); len(parts) > 0 {
			secretPaths = append(secretPaths, parts)
		}
	}

	if showEmptyGroups {
		for _, groupPath := range emptyGroupPaths {
			if parts := cleanPathParts(groupPath); len(parts) > 0 {
				groupPaths = append(groupPaths, parts)
			}
		}
	}

	allPaths := make([][]string, 0, len(secretPaths)+len(groupPaths))
	allPaths = append(allPaths, secretPaths...)
	allPaths = append(allPaths, groupPaths...)

	if len(allPaths) == 0 {
		fmt.Fprintln(os.Stderr, "No secrets found")
		return
	}

	if shouldStripTopGroup(allPaths) {
		for i := range secretPaths {
			secretPaths[i] = secretPaths[i][1:]
		}

		for i := range groupPaths {
			groupPaths[i] = groupPaths[i][1:]
		}
	}

	for _, parts := range secretPaths {
		addSecretPath(root, parts)
	}

	for _, parts := range groupPaths {
		addGroupPath(root, parts)
	}

	printTreeChildren(root, "")
}

func addSecretPath(root *treeNode, parts []string) {
	node := root

	for i, part := range parts {
		if node.children[part] == nil {
			node.children[part] = newTreeNode()
		}

		node = node.children[part]

		if i < len(parts)-1 {
			node.isGroup = true
		}
	}
}

func addGroupPath(root *treeNode, parts []string) {
	node := root

	for _, part := range parts {
		if node.children[part] == nil {
			node.children[part] = newTreeNode()
		}

		node = node.children[part]
		node.isGroup = true
	}
}

func cleanPathParts(p string) []string {
	parts := strings.Split(p, "/")
	clean := make([]string, 0, len(parts))

	for _, part := range parts {
		if part != "" {
			clean = append(clean, part)
		}
	}

	return clean
}

func shouldStripTopGroup(paths [][]string) bool {
	if len(paths) == 0 || len(paths[0]) < 2 {
		return false
	}

	first := paths[0][0]

	for _, parts := range paths {
		if len(parts) < 2 || parts[0] != first {
			return false
		}
	}

	return true
}

func printTreeChildren(node *treeNode, prefix string) {
	names := make([]string, 0, len(node.children))

	for name := range node.children {
		names = append(names, name)
	}

	sort.Strings(names)

	for i, name := range names {
		last := i == len(names)-1

		connector := "├── "
		nextPrefix := prefix + "│   "

		if last {
			connector = "└── "
			nextPrefix = prefix + "    "
		}

		child := node.children[name]
		displayName := name

		if child.isGroup {
			displayName += "/"
		}

		fmt.Println(prefix + connector + displayName)
		printTreeChildren(child, nextPrefix)
	}
}

func printHelp() {
	printLogoTitle()
	fmt.Println("Usage:")
	fmt.Println("  sato [flags] docker compose [compose-args...]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --db-path=PATH                            Path to KeePass-compatible .kdbx database")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help                                      Show this help")
	fmt.Println("  version                                   Show version")
	fmt.Println("  completion bash                           Show bash completion")
	fmt.Println()
	fmt.Println("  get secrets                               List secret names from KeePass-compatible .kdbx database")
	fmt.Println("  get secrets --show-empty-groups           List secret names and empty groups")
	fmt.Println("  get secrets --tree                        List secret names as a group tree")
	fmt.Println("  get secrets --tree --show-empty-groups    Include empty groups in tree output")
	fmt.Println("  get secret <NAME>                         Print one value of <NAME> secret to stdout")
	fmt.Println("  get secret <NAME> -q (--quiet)            Print only secret value, useful for scripts")
	fmt.Println()
	fmt.Println("  docker compose ...                        Run docker compose with KeePass-compatible secrets")

	fmt.Println()
	fmt.Printf("© %d Marcus Aprelius\n", time.Now().Year())
	fmt.Println("https://github.com/Marcus-Aprelius/sato")
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

func parseGetSecretArgs(args []string, inheritedQuiet bool) (string, bool, error) {
	secretName := ""
	quiet := inheritedQuiet

	for _, arg := range args {
		switch arg {
		case "-q", "--quiet":
			quiet = true
		default:
			if strings.HasPrefix(arg, "-") {
				return "", quiet, fmt.Errorf("unknown option: %s", arg)
			}

			if secretName != "" {
				return "", quiet, errors.New("usage: sato get secret NAME [-q|--quiet]")
			}

			secretName = arg
		}
	}

	if secretName == "" {
		return "", quiet, errors.New("usage: sato get secret NAME [-q|--quiet]")
	}

	return secretName, quiet, nil
}

func findSecretValue(entries []SecretEntry, name string) (string, error) {
	matches := make([]SecretEntry, 0)

	for _, entry := range entries {
		if entry.Name == name {
			matches = append(matches, entry)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("secret not found: %s", name)
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("secret name is not unique: %s", name)
	}

	return matches[0].Value, nil
}

func Run() error {
	dbPathFlag := flag.String("db-path", "", "path to KeePass-compatible .kdbx database")
	qFlag := flag.Bool("q", false, "print only secret value for get secret")
	quietFlag := flag.Bool("quiet", false, "print only secret value for get secret")
	flag.Parse()

	quietMode := *qFlag || *quietFlag

	args := flag.Args()

	if len(args) > 0 {
		switch args[0] {
		case "version":
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

			os.Exit(0)

		case "completion":
			shell := "bash"
			if len(args) > 1 {
				shell = args[1]
			}
			RunCompletion(shell)
			os.Exit(0)

		case "help":
			printHelp()
			os.Exit(0)

		case "get":
			if len(args) > 1 && args[1] == "secret" {
				secretName, commandQuiet, err := parseGetSecretArgs(args[2:], quietMode)
				if err != nil {
					return err
				}

				dbPath := FindDBPath(*dbPathFlag)
				if dbPath == "" {
					return errors.New("DB not found")
				}

				if !commandQuiet {
					fmt.Fprintf(os.Stderr, "[DB: %s]\n", dbPath)
				}

				password, err := ReadPasswordWithPrompt(!commandQuiet)
				if err != nil {
					return err
				}

				data, err := LoadKeePassData(dbPath, password)
				if err != nil {
					return err
				}

				value, err := findSecretValue(data.Entries, secretName)
				if err != nil {
					return err
				}

				fmt.Println(value)
				os.Exit(0)
			}
			if len(args) > 1 && args[1] == "secrets" {
				treeMode := false
				showEmptyGroups := false

				for _, arg := range args[2:] {
					switch arg {
					case "--tree":
						treeMode = true
					case "--show-empty-groups":
						showEmptyGroups = true
					default:
						return fmt.Errorf("unknown option: %s", arg)
					}
				}

				dbPath := FindDBPath(*dbPathFlag)
				if dbPath == "" {
					return errors.New("DB not found")
				}

				fmt.Fprintf(os.Stderr, "[DB: %s]\n", dbPath)

				password, err := ReadPassword()
				if err != nil {
					return err
				}

				data, err := LoadKeePassData(dbPath, password)
				if err != nil {
					return err
				}

				if treeMode {
					printSecretsTree(data.Entries, data.EmptyGroupPaths, showEmptyGroups)
					os.Exit(0)
				}

				values := make([]string, 0, len(data.Entries)+len(data.EmptyGroupPaths))

				for _, entry := range data.Entries {
					values = append(values, entry.Name)
				}

				if showEmptyGroups {
					values = append(values, emptyGroupDisplayPaths(data.Entries, data.EmptyGroupPaths)...)
				}

				sort.Strings(values)

				for _, value := range values {
					fmt.Println(value)
				}

				os.Exit(0)
			}
		}
	}

	if len(args) == 0 {
		printStatus(*dbPathFlag)
		os.Exit(0)
	}

	if !IsAllowedCommand(args) {
		fmt.Fprintf(os.Stderr, "ERROR: Only 'docker compose' is supported\n")
		fmt.Fprintf(os.Stderr, "\n")
		printStatus(*dbPathFlag)
		os.Exit(1)
	}

	if !HasComposeSubcommand(args) {
		fmt.Fprintf(os.Stderr, "ERROR: docker compose subcommand is required\n")
		fmt.Fprintf(os.Stderr, "\n")
		printStatus(*dbPathFlag)
		os.Exit(1)
	}

	dbPath := FindDBPath(*dbPathFlag)

	if dbPath == "" {
		fmt.Fprintf(os.Stderr, "ERROR: File .kdbx not found in configured locations\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "[DB: %s]\n", dbPath)

	password, err := ReadPassword()
	if err != nil {
		return err
	}

	env, err := LoadKeePass(dbPath, password)
	if err != nil {
		return err
	}

	return RunCommand(args[0], args[1:], env)
}
