package main

import (
	"log"
)

type installCommand struct{}

func (cmd *installCommand) Execute(args []string) error {
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

			if err := installSelf(installPath); err != nil {
				return appError{err, "Could not install locally. Depending on your system's configuration, you may need to run the install as an administrator."}
			}

			return nil
		}

		token := parseNameVersionToken(args[0])
		gamePath := getGameOrDefault(current.config.Games, "")
		mod, err := installMod(current.config, gamePath, current.profile, token)
		if err != nil {
			return err
		}

		log.Printf("%s@%s installed at %s\n", mod.Name, mod.SemVersion, mod.source)
		return nil
	})
}

func init() {
	install, _ := commandParser.AddCommand("install", "", "", &installCommand{})
	install.Aliases = []string{"i"}
}
