package main

import (
	//"fmt"
	"log"
)

type uninstallCommand struct{}

func (cmd *uninstallCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext) error {
		self := args[0] == "self"
		if self {
			installPath := ""
			if len(args) > 1 {
				installPath = args[1]
			}

			if err := uninstallSelf(installPath); err != nil {
				return appError{err, "Could not uninstall. Depending on your system's configuration, you may need to run the uninstall again as an administrator."}
			}

			return nil
		}

		name := args[0]
		gamePath := getGameOrDefault(current.config.Games, name)
		mod, err := uninstallMod(current.config, gamePath, current.profile, name)
		if err != nil {
			return appError{err, "Could not uninstall mod"}
		}

		log.Printf("%s@%s uninstalled\n", mod.Name, mod.SemVersion)
		return nil
	})
}

func init() {
	c, _ := commandParser.AddCommand("uninstall", "uninstall a mod", "Uninstall a mod from the current profile.", &uninstallCommand{})
	c.Aliases = []string{"u"}
}
