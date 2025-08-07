package klink

import (
	"fmt"
	"os"
	"strings"
)

type CliFunction func([]string, map[string]string)

type CliNode struct {
	name        string
	description string
	options     []*CliNode
	action      CliFunction
	flag        []CliFlag
}

type CliFlag struct {
	name        string
	description string
}

func GetActions() []string {
	args := os.Args[1:]
	actions := []string{}
	for k, v := range args {
		if strings.HasPrefix(v, "--") {
			continue
		}

		if strings.HasPrefix(v, "-") {
			continue
		}

		if k != 0 && strings.HasPrefix(args[k-1], "-") {
			continue
		}

		if k != 0 && strings.HasPrefix(args[k-1], "--") {
			continue
		}

		actions = append(actions, v)
	}
	return actions
}

func GetFlags() map[string]string {
	args := os.Args
	flags := map[string]string{}

	for k, v := range args {
		if strings.HasPrefix(v, "--") {
			//  TODO: Handle set like --flag=value

			if len(args) > k+1 {
				flags[v] = args[k+1]
			} else {
				flags[v] = "true"
			}
			continue
		}

		if strings.HasPrefix(v, "-") {
			continue
		}
	}
	return flags
}

func SubNode(parent *CliNode, node *CliNode) *CliNode {
	parent.options = append(parent.options, node)
	return node
}

func PrintHelp(node *CliNode) {
	fmt.Printf("%s\n\n", node.description)

	fmt.Println("Available commands:")
	for _, v := range node.options {
		fmt.Printf("  %-10s %-24s\n", v.name, v.description)
	}
	fmt.Println("")

	fmt.Println("Flags:")
	for _, v := range node.flag {
		fmt.Printf("  %-10s %-24s\n", v.name, v.description)
	}
}

func processNode(node *CliNode, actions []string, flags map[string]string) {
	// If has action, go no lower
	if node.action != nil {
		// fmt.Printf("\tProcessing action for node %s\n", node.name)
		node.action(actions, flags)
		return
	}

	// If no more actions, don't even bother looping
	if len(actions) == 0 {
		PrintHelp(node)
		return
	}

	action := actions[0]
	for _, v := range node.options {
		if v.name == action {
			processNode(v, actions[1:], flags)
			return
		}
	}

	// If all else failed, just print the help for the current node
	PrintHelp(node)
}
