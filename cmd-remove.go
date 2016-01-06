package main

import (
	"fmt"
	//"os"
)

/* Remove Command */
type removeCommand struct{}

func (cmd *removeCommand) Usage() string {
	return ""
}

func (cmd *removeCommand) Execute(args []string) error {
	return usageError{}
}

/* Remove Profile Command */
type removeProfileCommand struct{}

func (cmd *removeProfileCommand) Usage() string {
	return "<profile>"
}

func (cmd *removeProfileCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext, log logger) error {
		profileId := current.args[0]
		profile, ok := current.profiles[profileId]
		if !ok {
			return appError{fmt.Sprintf("Could not find the profile \"%s\"", profileId)}
		}

		err := dropProfile(profile)
		if err != nil {
			return err
		}

		log("\"%s\" has been removed\n", profile.Name)
		return nil
	})
}

/* Remove Game Command */
type removeGameCommand struct{}

func (cmd *removeGameCommand) Usage() string {
	return "<game>"
}

func (cmd *removeGameCommand) Execute(args []string) error {
	if len(args) < 0 {
		return usageError{}
	}

	return runInContext(func(current *executionContext, log logger) error {
		alias := args[0]
		gamePath, ok := current.config.Games[alias]
		if !ok {
			// TODO: embed error
			return appError{fmt.Sprintf("game \"%s\" does not exist.\n", alias)}
		}

		err := detachGameFolder(current.config, alias)
		if err != nil {
			// TODO: embed error
			return appError{"Could not remove game from manager"}
		}

		err = saveManagerConfig(current.config, current.homePath)
		if err != nil {
			// TODO: embed error
			return appError{"Could not save config"}
		}

		log("game removed: %s => %s\n", alias, gamePath)
		return nil
	})
}

/* Remove Channel Command */
type removeChannelCommand struct{}

func (cmd *removeChannelCommand) Usage() string {
	return ""
}

func (cmd *removeChannelCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext, log logger) error {
		alias := args[0]
		var channel string
		ok := false
		if alias == "default" && current.config.IncludeDefaultChannel {
			channel = getDefaultRemote()
			ok = true
		} else {
			channel, ok = current.config.Remotes[alias]
		}

		if !ok {
			return appError{"channel not found"}
		}

		err := removeChannel(alias, current.config)
		if err != nil {
			// TODO: embed error
			return appError{"couldn't remove channel"}
		}

		err = saveManagerConfig(current.config, current.homePath)
		if err != nil {
			return appError{"couldn't save manager config"}
		}

		log("channel removed: %s => %s\n", alias, channel)
		return nil
	})
}

func init() {
	rm, _ := commandParser.AddCommand("remove", "remove manager assets", "Removes a manager asset", &removeCommand{})
	rm.AddCommand("profile", "remove a profile", "Removes the profile from the manager", &removeProfileCommand{})
	rm.AddCommand("game", "remove a game folder", "Removes a game by the specified alias from the manager.", &removeGameCommand{})
	rm.AddCommand("channel", "remove a channel", "Removes a channel from the manager", &removeChannelCommand{})
}
