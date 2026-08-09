package sato

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const bashCompletion = `# sato bash completion
_sato_docker_completion() {
	local old_words old_cword

	if ! declare -F _docker >/dev/null 2>&1; then
		return 1
	fi

	old_words=("${COMP_WORDS[@]}")
	old_cword="$COMP_CWORD"

	COMP_WORDS=("docker" "${old_words[@]:2}")
	COMP_CWORD=$((old_cword - 1))

	_docker

	COMP_WORDS=("${old_words[@]}")
	COMP_CWORD="$old_cword"

	return 0
}

_sato_completion() {
	local cur prev words cword
	_init_completion || return

	local commands="docker version completion help get git"
	local flags="--db-path --secret"
	local docker_compose_opts="attach build config create events export kill ls port publish push rm scale stats top up volumes watch bridge commit cp down exec images logs pause ps pull restart run start stop unpause version wait"

	if [[ "${words[1]}" == "docker" ]]; then
		if _sato_docker_completion; then
			return
		fi

		case "$prev" in
			docker)
				COMPREPLY=( $(compgen -W "compose help" -- "$cur") )
				return
				;;
			compose)
				COMPREPLY=( $(compgen -W "${docker_compose_opts}" -- "$cur") )
				return
				;;
		esac

		case "$cword" in
			2)
				COMPREPLY=( $(compgen -W "compose help" -- "$cur") )
				return
				;;
			3)
				if [[ "${words[2]}" == "compose" ]]; then
					COMPREPLY=( $(compgen -W "${docker_compose_opts}" -- "$cur") )
					return
				fi
				;;
		esac

		return
	fi

	case "$prev" in
		get)
			COMPREPLY=( $(compgen -W "secret secrets" -- "$cur") )
			return
			;;
		git)
			COMPREPLY=( $(compgen -W "clone fetch pull push help" -- "$cur") )
			return
			;;
		completion)
			COMPREPLY=( $(compgen -W "bash" -- "$cur") )
			return
			;;
		bash)
			if [[ "${words[1]}" == "completion" ]]; then
				COMPREPLY=( $(compgen -W "add delete update status" -- "$cur") )
				return
			fi
			;;
		--db-path)
			_filedir '*.kdbx'
			return
			;;
		--secret)
			COMPREPLY=()
			return
			;;
	esac

	case "$cword" in
		1)
			COMPREPLY=( $(compgen -W "${commands} ${flags}" -- "$cur") )
			;;
		2)
			if [[ "${words[1]}" == "--db-path" ]]; then
				_filedir '*.kdbx'
			elif [[ "${words[1]}" == "--secret" ]]; then
				COMPREPLY=()
			elif [[ "${words[1]}" == "get" ]]; then
				COMPREPLY=( $(compgen -W "secret secrets" -- "$cur") )
			elif [[ "${words[1]}" == "git" ]]; then
				COMPREPLY=( $(compgen -W "clone fetch pull push help" -- "$cur") )
			elif [[ "${words[1]}" == "docker" ]]; then
				COMPREPLY=( $(compgen -W "compose help" -- "$cur") )
			elif [[ "${words[1]}" == "completion" ]]; then
				COMPREPLY=( $(compgen -W "bash" -- "$cur") )
			fi
			;;
		3)
			if [[ "${words[1]}" == "get" && "${words[2]}" == "secrets" ]]; then
				COMPREPLY=( $(compgen -W "--tree --show-empty-groups" -- "$cur") )
			elif [[ "${words[1]}" == "get" && "${words[2]}" == "secret" ]]; then
				COMPREPLY=()
			elif [[ "${words[1]}" == "completion" && "${words[2]}" == "bash" ]]; then
				COMPREPLY=( $(compgen -W "add delete update status" -- "$cur") )
			fi
			;;
		4)
			if [[ "${words[1]}" == "get" && "${words[2]}" == "secrets" ]]; then
				case " ${words[*]} " in
					*" --tree "*)
						COMPREPLY=( $(compgen -W "--show-empty-groups" -- "$cur") )
						;;
					*" --show-empty-groups "*)
						COMPREPLY=( $(compgen -W "--tree" -- "$cur") )
						;;
					*)
						COMPREPLY=( $(compgen -W "--tree --show-empty-groups" -- "$cur") )
						;;
				esac
			elif [[ "${words[1]}" == "get" && "${words[2]}" == "secret" ]]; then
				COMPREPLY=( $(compgen -W "-q --quiet" -- "$cur") )
			fi
			;;
	esac
}

complete -F _sato_completion sato
`

const bashCompletionStartMarker = "# sato completion start"
const bashCompletionEndMarker = "# sato completion end"

const bashCompletionBlock = `
# sato completion start
source ~/.sato-completion.bash
# sato completion end
`

func RunCompletion(shell string, action ...string) {
	if shell != "bash" {
		fmt.Fprintf(os.Stderr, "ERROR: unsupported shell %q. Supported: bash\n", shell)
		os.Exit(1)
	}

	if len(action) == 0 {
		fmt.Print(bashCompletion)
		return
	}

	switch action[0] {
	case "add":
		if err := addBashCompletion(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("SATO bash completion: Added")
		fmt.Println("Reload shell:")
		fmt.Println("source ~/.bashrc")

	case "delete":
		if err := deleteBashCompletion(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("SATO bash completion: Removed")
		fmt.Println("Reload shell:")
		fmt.Println("source ~/.bashrc")

	case "update":
		if err := updateBashCompletion(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("SATO bash completion: Updated")
		fmt.Println("Reload shell:")
		fmt.Println("source ~/.bashrc")

	case "status":
		enabled, err := bashCompletionEnabled()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		if enabled {
			fmt.Println("SATO bash completion: Enabled")
		} else {
			fmt.Println("SATO bash completion: Disabled")
		}

	default:
		fmt.Fprintf(
			os.Stderr,
			"ERROR: unknown completion action %q. Supported: add, delete, update, status\n",
			action[0],
		)
		os.Exit(1)
	}
}

func bashrcPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".bashrc"), nil
}

func completionFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".sato-completion.bash"), nil
}

func writeCompletionFile() error {
	path, err := completionFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(bashCompletion), 0o644)
}

func addBashCompletion() error {
	if err := writeCompletionFile(); err != nil {
		return err
	}

	path, err := bashrcPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	content := string(data)
	if strings.Contains(content, bashCompletionStartMarker) {
		return nil
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(bashCompletionBlock)
	return err
}

func deleteBashCompletion() error {
	bashrc, err := bashrcPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(bashrc)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil {
		content := removeBashCompletionBlock(string(data))
		if err := os.WriteFile(bashrc, []byte(content), 0o644); err != nil {
			return err
		}
	}

	completionFile, err := completionFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(completionFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func updateBashCompletion() error {
	enabled, err := bashCompletionEnabled()
	if err != nil {
		return err
	}

	if err := writeCompletionFile(); err != nil {
		return err
	}

	if !enabled {
		return addBashCompletion()
	}

	return nil
}

func bashCompletionEnabled() (bool, error) {
	bashrc, err := bashrcPath()
	if err != nil {
		return false, err
	}

	data, err := os.ReadFile(bashrc)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	completionFile, err := completionFilePath()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(completionFile); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return strings.Contains(string(data), bashCompletionStartMarker), nil
}

func removeBashCompletionBlock(content string) string {
	start := strings.Index(content, bashCompletionStartMarker)
	if start < 0 {
		return content
	}

	end := strings.Index(content[start:], bashCompletionEndMarker)
	if end < 0 {
		return content
	}

	end += start + len(bashCompletionEndMarker)

	if end < len(content) && content[end] == '\n' {
		end++
	}

	before := content[:start]
	after := content[end:]

	return strings.TrimRight(before, "\n") + "\n" + strings.TrimLeft(after, "\n")
}
