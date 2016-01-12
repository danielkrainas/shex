package main

import (
	"log"
)

type useCommand struct{}

func (cmd *useCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext) error {
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

			log.Printf("active profile set to: %s\n", newProfile.Name)
		} else {
			log.Printf("profile already active")
		}

		return nil
	})
}

func init() {
	commandParser.AddCommand("use", "", "", &useCommand{})
}
