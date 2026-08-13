package sato

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
)

func Run() error {
	dbPathFlag := flag.String("db-path", "", "path to KeePass-compatible .kdbx database")
	secretFlag := flag.String("secret", "", "KeePass secret name override")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		switch args[0] {

		case "version":
			RunVersion()
			os.Exit(0)

		case "completion":
			shell := "bash"
			if len(args) > 1 {
				shell = args[1]
			}

			if len(args) > 2 {
				RunCompletion(shell, args[2:]...)
			} else {
				RunCompletion(shell)
			}

			os.Exit(0)

			RunCompletion(shell)
			os.Exit(0)

		case "help":
			printHelp()
			os.Exit(0)

		case "docker":
			if len(args) == 1 {
				printDockerHelp()
				os.Exit(0)
			}

			if args[1] == "help" {
				printDockerHelp()
				os.Exit(0)
			}

		case "git":
			if len(args) == 1 {
				printGitHelp()
				os.Exit(0)
			}

			if args[1] == "help" {
				printGitHelp()
				os.Exit(0)
			}

		case "get":
			if len(args) > 1 && args[1] == "secret" {
				secretName, commandQuiet, err := parseGetSecretArgs(args[2:], false)
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

	if IsAllowedGitCommand(args) {
		secretName, gitArgs, err := gitSecretFromArgs(args, *secretFlag)
		if err != nil {
			return err
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

		token, err := findSecretValue(data.Entries, secretName)
		if err != nil {
			return err
		}

		return RunGitCommand(gitArgs, token)
	}

	if !IsAllowedCommand(args) {
		fmt.Fprintf(os.Stderr, "ERROR: Only 'docker compose' and 'git clone|fetch|pull|push' are supported\n")
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
