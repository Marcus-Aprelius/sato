package sato

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

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
