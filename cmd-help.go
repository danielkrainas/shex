package main

import (
	"fmt"
	"os"
)

type helpCommand struct{}

func (cmd *helpCommand) Execute(args []string) error {
	showGeneralHelp := commandParser.Active == nil || commandParser.Active == commandParser.Command || len(args) < 1 && commandParser.Active.Name == "help"
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
	} else if len(args) > 0 {
		c := commandParser.Find(args[0])
		commandParser.Active = c
		commandParser.WriteHelp(os.Stdout)
	} else {
		println(args)
	}

	return nil
}

func init() {
	commandParser.AddCommand("help", "displays help information", "Display help and usage information.", &helpCommand{})
}
