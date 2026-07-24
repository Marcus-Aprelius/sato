package sato

import (
	"fmt"
	"os"
)

const bashCompletion = `# sato bash completion
_sato_docker_completion() {
	local old_words old_cword

	# Docker bash completion must be loaded
	if ! declare -F _docker >/dev/null 2>&1; then
		return 1
	fi

	old_words=("${COMP_WORDS[@]}")
	old_cword="$COMP_CWORD"

	# Rewrite:
	#   sato docker compose up
	# to:
	#   docker compose up
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

	local commands="docker version completion help get"
	local flags="--db-path"
	local docker_compose_opts="attach build config create events export kill ls port publish push rm scale stats top up volumes watch bridge commit cp down exec images logs pause ps pull restart run start stop unpause version wait"

	# Delegate docker compose completion to Docker completion if available
	if [[ "${words[1]}" == "docker" ]]; then
		if _sato_docker_completion; then
			return
		fi

		# Fallback if Docker completion is not loaded
		case "$prev" in
			docker)
				COMPREPLY=( $(compgen -W "compose" -- "$cur") )
				return
				;;
			compose)
				COMPREPLY=( $(compgen -W "${docker_compose_opts}" -- "$cur") )
				return
				;;
		esac

		case "$cword" in
			2)
				COMPREPLY=( $(compgen -W "compose" -- "$cur") )
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
			COMPREPLY=( $(compgen -W "secrets" -- "$cur") )
			return
			;;
		--db-path)
			_filedir '*.kdbx'
			return
			;;
		completion)
			COMPREPLY=( $(compgen -W "bash" -- "$cur") )
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
			elif [[ "${words[1]}" == "get" ]]; then
				COMPREPLY=( $(compgen -W "secrets" -- "$cur") )
			elif [[ "${words[1]}" == "completion" ]]; then
				COMPREPLY=( $(compgen -W "bash" -- "$cur") )
			fi
			;;
		3)
			if [[ "${words[1]}" == "get" && "${words[2]}" == "secrets" ]]; then
				COMPREPLY=( $(compgen -W "--tree" -- "$cur") )
			fi
			;;
		4)
			if [[ "${words[1]}" == "get" && "${words[2]}" == "secrets" && "${words[3]}" == "--tree" ]]; then
				COMPREPLY=( $(compgen -W "--show-empty-groups" -- "$cur") )
			fi
			;;
	esac
}

complete -F _sato_completion sato
`

// Generate and print the completion script for the shell
func RunCompletion(shell string) {
	switch shell {
	case "bash":
		fmt.Print(bashCompletion)
	default:
		fmt.Fprintf(os.Stderr, "ERROR: unsupported shell %q. Supported: bash\n", shell)
		os.Exit(1)
	}
}
