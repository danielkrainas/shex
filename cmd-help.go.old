// -build

package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/kardianos/osext"
	"strings"
)

type helpCommand struct{}

func (cmd *helpCommand) Execute(args []string) error {
	termKey := strings.Join(args, " ")
	showGeneralHelp := commandParser.Active == nil || commandParser.Active == commandParser.Command || len(args) < 1 && commandParser.Active.Name == "help"
	if !showGeneralHelp && len(args) > 0 {
		c := commandParser.Find(args[0])
		help := helpIndex[termKey]
		var p *flags.Command
		if c != nil && len(args) > 1 {
			p = c
			c = c.Find(args[1])
		}

		if help == nil && c == nil {
			fmt.Printf("no help for: \"%s\"\n", termKey)
			return nil
		} else if help != nil {
			c = p
			p = nil
		}

		command := termKey
		var subs []*flags.Command
		if c != nil {
			subs = c.Commands()
			if len(subs) > 0 {
				command += " <command>"
			}
		}

		fmt.Printf("\n%s\n\n", help.long)
		fmt.Printf("Usage:\n    %s %s %s\n\n", commandParser.Command.Name, command, help.usage)

		if subs != nil && len(subs) > 0 {
			fmt.Printf("where <command> is one of:\n")
			rowCount := 0
			for i, c := range subs {
				if rowCount == 0 {
					fmt.Print("    ")
				}

				fmt.Print(c.Name)
				if (i + 1) < len(subs) {
					fmt.Print(", ")
				}

				rowCount++
				if rowCount >= 6 {
					rowCount = 0
					fmt.Print("\n")
				}
			}

			if len(subs) < 6 {
				fmt.Print("\n")
			}

			fmt.Printf("\nUse \"%s help %s <cmd>\" to display help about <cmd>\n\n", commandParser.Command.Name, c.Name)
		}
	} else {
		showGeneralHelp = true
	}

	if showGeneralHelp {
		_ = runInContext(func(current *executionContext) error {
			binaryPath, _ := osext.ExecutableFolder()
			fmt.Printf("\nUsage:\n    %s <command> [options]\n\n", commandParser.Command.Name)

			fmt.Print("where <command> is one of:\n")
			rowCount := 0
			for i, c := range commandParser.Commands() {
				if rowCount == 0 {
					fmt.Print("    ")
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

			fmt.Printf("\n\nUse \"%s help <cmd>\" to display help about <cmd>\n\n", commandParser.Command.Name)
			fmt.Printf("Config:\n    %s\n", current.config.filePath)
			fmt.Printf("Home:\n    %s\n", current.homePath)
			fmt.Printf("\n%s@%s %s\n", commandParser.Command.Name, buildVersion, binaryPath)
			return nil
		})
	}

	return nil
}

func init() {
	commandParser.AddCommand("help", "", "", &helpCommand{})
}
