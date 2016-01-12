package main

import (
	"fmt"
	"log"
)

/* Remove Command */
type removeCommand struct{}

func (cmd *removeCommand) Execute(args []string) error {
	return usageError{}
}

/* Remove Profile Command */
type removeProfileCommand struct{}

func (cmd *removeProfileCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext) error {
		profileId := args[0]
		profile, ok := current.profiles[profileId]
		if !ok {
			return appError{nil, fmt.Sprintf("Could not find the profile \"%s\"", profileId)}
		}

		err := dropProfile(profile)
		if err != nil {
			return err
		}

		log.Printf("\"%s\" has been removed\n", profile.Name)
		return nil
	})
}

/* Remove Game Command */
type removeGameCommand struct{}

func (cmd *removeGameCommand) Execute(args []string) error {
	if len(args) < 0 {
		return usageError{}
	}

	return runInContext(func(current *executionContext) error {
		alias := args[0]
		gamePath, ok := current.config.Games[alias]
		if !ok {
			return appError{nil, fmt.Sprintf("game \"%s\" does not exist.\n", alias)}
		}

		err := detachGameFolder(current.config, alias)
		if err != nil {
			return appError{err, "Could not remove game from manager"}
		}

		err = saveManagerConfig(current.config, current.homePath)
		if err != nil {
			return appError{err, "Could not save config"}
		}

		log.Printf("game removed: %s => %s\n", alias, gamePath)
		return nil
	})
}

/* Remove Channel Command */
type removeChannelCommand struct{}

func (cmd *removeChannelCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext) error {
		alias := args[0]
		var channel *Channel
		ok := false
		if alias == "default" && current.config.IncludeDefaultChannel {
			channel = defaultChannel
			ok = true
		} else {
			channel, ok = current.channels[alias]
		}

		if !ok {
			return appError{nil, "channel not found"}
		}

		if channel == defaultChannel {
			current.config.IncludeDefaultChannel = false
			if err := saveManagerConfig(current.config, current.homePath); err != nil {
				return appError{err, "couldn't save channel"}
			}
		} else if err := channel.remove(); err != nil {
			return appError{err, "couldn't remove channel"}
		}

		log.Printf("channel removed: %s => %s\n", channel.Alias, channel.Endpoint)
		return nil
	})
}

func init() {
	rm, _ := commandParser.AddCommand("remove", "", "", &removeCommand{})
	rm.AddCommand("profile", "", "", &removeProfileCommand{})
	rm.AddCommand("game", "", "", &removeGameCommand{})
	rm.AddCommand("channel", "", "", &removeChannelCommand{})
}
