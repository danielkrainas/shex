package main

type useCommand struct{}

func (cmd *useCommand) Execute(args []string) error {
	// TODO check arg length and show usage

	return runInContext(func(current *executionContext, log logger) error {
		// TODO: this needs to uninstall existing mods and install mods from new profile
		profiles := current.profiles
		config := current.config

		newProfileName := args[0]
		if newProfileName != config.ActiveProfile {
			newProfile := profiles[newProfileName]
			config.ActiveProfile = newProfile.Id
			err := saveManagerConfig(config, current.homePath)
			if err != nil {
				return err
			}

			log("active profile set to: %s\n", newProfile.Name)
		} else {
			log("profile already active")
		}

		return nil
	})
}

func (cmd *useCommand) Usage() string {
	return "use <profile>"
}

func init() {
	commandParser.AddCommand("use", "sets the active profile", "sets the active mod profile", &useCommand{})
}
