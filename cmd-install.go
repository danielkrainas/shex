package main

import ()

type installCommand struct{}

func (cmd *installCommand) Usage() string {
	return "self|<mod>"
}

func (cmd *installCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext, log logger) error {
		self := args[0] == "self"
		if self {
			installPath := ""
			if len(args) > 1 {
				installPath = args[1]
			}

			if err := installSelf(installPath); err != nil {
				println(err.Error())
				return appError{"Could not install locally. Depending on your system's configuration, you may need to run the install as an administrator."}
			}

			return nil
		}

		token := parseNameVersionToken(args[0])
		gamePath := getGameOrDefault(current.config.Games, "")
		mod, err := installMod(current.config, gamePath, current.profile, token)
		if err != nil {
			return err
		}

		log("%s@%s installed at %s\n", mod.Name, mod.SemVersion, mod.source)
		return nil
	})
}

func init() {
	install, _ := commandParser.AddCommand("install", "install a mod", "Install a mod into the current profile", &installCommand{})
	install.Aliases = []string{"i"}
	install.ArgsRequired = true
}
