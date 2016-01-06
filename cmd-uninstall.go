package main

import (
	"fmt"
)

type uninstallCommand struct{}

func (cmd *uninstallCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext, log logger) error {
		self := args[0] == "self"
		if self {
			installPath := ""
			if len(args) > 0 {
				installPath = args[0]
			}

			if err := uninstallSelf(installPath); err != nil {
				fmt.Printf(err.Error())
				return appError{"Could not uninstall. Depending on your system's configuration, you may need to run the uninstall again as an administrator."}
			}

			return nil
		}

		name := args[0]
		gamePath := getGameOrDefault(current.config.Games, name)
		mod, err := uninstallMod(current.config, gamePath, current.profile, name)
		if err != nil {
			// TODO: embed error
			return appError{"Could not uninstall mod"}
		}

		log("%s@%s uninstalled\n", mod.Name, mod.SemVersion)
		return nil
	})
}

func init() {
	c, _ := commandParser.AddCommand("uninstall", "uninstall a mod", "Uninstall a mod from the current profile.", &uninstallCommand{})
	c.Aliases = []string{"u"}
}
