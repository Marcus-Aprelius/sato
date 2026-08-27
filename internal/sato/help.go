package sato

import (
	"fmt"
	"time"
)

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

func printHelp() {
	printLogoTitle()
	fmt.Println("Usage:")
	fmt.Println("  sato [flags] [commands]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --db-path=<PATH>                          Path to a Secret database (.kdbx, .psafe3, .ibak)")
	fmt.Println("  --secret=<NAME>                           Override secret name for Git commands")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  help                                      Show this help")
	fmt.Println("  git help                                  Show git command help")
	fmt.Println("  docker help                               Show docker command help")
	fmt.Println()
	fmt.Println("  get secrets                               List secret names from Secret database (.kdbx, .psafe3, .ibak)")
	fmt.Println("  get secrets --show-empty-groups           List secret names and empty groups")
	fmt.Println("  get secrets --tree                        List secret names as a group tree")
	fmt.Println("  get secrets --tree --show-empty-groups    Include empty groups in tree output")
	fmt.Println()
	fmt.Println("  get secret <NAME>                         Print one secret value to stdout")
	fmt.Println("  get secret <NAME> [-q|--quiet]            Print only secret value, useful for scripts")
	fmt.Println()
	fmt.Println("  version                                   Show version")
	fmt.Println("  completion bash                           Show bash completion code")
	fmt.Println("  completion bash status                    Show bash completion status")
	fmt.Println("  completion bash add|delete|update         Add, delete, or update bash completion in ~/.bashrc")
	fmt.Println()
	fmt.Printf("© %d Marcus Aprelius\n", time.Now().Year())
	fmt.Println("https://github.com/Marcus-Aprelius/sato")
}

func printDockerHelp() {
	printLogoTitle()
	fmt.Println()
	fmt.Println("Reads secrets from supported databases and injects them into Docker Compose commands.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  sato [flags] [commands]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --db-path=<PATH>                          Path to a secret database (.kdbx, .psafe3, .ibak)")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  docker compose ...                        Run Docker Compose with secrets from the configured database")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  sato --db-path=~/secrets.psafe3 docker compose config")
	fmt.Println("  sato docker compose up -d")
	fmt.Println()
	fmt.Printf("© %d Marcus Aprelius\n", time.Now().Year())
	fmt.Println("https://github.com/Marcus-Aprelius/sato")
}

func printGitHelp() {
	printLogoTitle()
	fmt.Println()
	fmt.Println("Uses the repository name as the default secret name for Git commands.")
	fmt.Println("  To override the secret name, use `--secret <NAME>` flag.")
	fmt.Println("Secret lookup:")
	fmt.Println("  clone:           SATO uses repository name from URL as secret name.")
	fmt.Println("  fetch/pull/push: SATO uses repository name from git remote origin URL.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  sato [flags] [commands]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --db-path=<PATH>                          Path to a Secret database (.kdbx, .psafe3, .ibak)")
	fmt.Println("  --secret=<NAME>                           Name of the secret for `sato git ...`")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  git clone <URL>                           Clone an HTTPS repository using a database secret")
	fmt.Println("  git fetch                                 Fetch using a database secret for the origin repository")
	fmt.Println("  git pull                                  Pull using a database secret for the origin repository")
	fmt.Println("  git push                                  Push using a database secret for the origin repository")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  sato --secret=sato-token git clone https://github.com/Marcus-Aprelius/sato.git")
	fmt.Println("  sato --db-path=~/secrets.psafe3 git clone https://github.com/Marcus-Aprelius/sato.git")
	fmt.Println()
	fmt.Printf("© %d Marcus Aprelius\n", time.Now().Year())
	fmt.Println("https://github.com/Marcus-Aprelius/sato")
}
