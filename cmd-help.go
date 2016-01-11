package main

import (
	"fmt"
	"os"
	"strings"
	//"text/template"
)

type helpCommand struct{}

func (cmd *helpCommand) Execute(args []string) error {
	showGeneralHelp := commandParser.Active == nil || commandParser.Active == commandParser.Command || len(args) < 1 && commandParser.Active.Name == "help"
	if !showGeneralHelp && len(args) > 0 {
		c := commandParser.Find(args[0])
		if c != nil && len(args) > 1 {
			c = c.Find(args[1])
		} else if c == nil {
			fmt.Printf("no help for: \"%s\"\n", strings.Join(args, " "))
			return nil
		}

		commandParser.Active = c
		commandParser.WriteHelp(os.Stdout)
	} else {
		showGeneralHelp = true
	}

	if showGeneralHelp {
		fmt.Printf("Usage:\n  %s <command> [options]\n\n", commandParser.Command.Name)

		fmt.Print("where <command> is one of:\n")
		rowCount := 0
		for i, c := range commandParser.Commands() {
			if rowCount == 0 {
				fmt.Print("  ")
			}

			fmt.Printf(c.Name)
			if (i + 1) < len(commandParser.Commands()) {
				fmt.Print(", ")
			}

			rowCount++
			if rowCount >= 6 {
				rowCount = 0
				fmt.Print("\n")
			}
		}

		fmt.Println("")
	}

	return nil
}

func init() {
	commandParser.AddCommand("help", "displays help information", "Display help and usage information.", &helpCommand{})
}
