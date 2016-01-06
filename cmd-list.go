package main

import (
	"fmt"
)

/* List Command */
type listCommand struct{}

func (cmd *listCommand) Execute(args []string) error {
	return usageError{}
}

/* List Mods Command */
type listModsCommand struct {
	Profile string `short:"p" long:"profile" description:"display mods installed in a profile"`
}

func (cmd *listModsCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext, log logger) error {
		profileOption := commandParser.FindOptionByLongName("profile")
		useProfile := profileOption != nil && profileOption.IsSet()
		profileName := ""
		var mods ModList
		if useProfile {
			profileName = cmd.Profile
			if len(cmd.Profile) > 0 {
				selectedProfile, ok := current.profiles[profileName]
				if !ok {
					return appError{fmt.Sprintf("profile %s not found", profileName)}
				}

				mods = selectedProfile.Mods
			} else {
				profileName = current.profile.Name
				mods = current.profile.Mods
			}
		} else if len(current.config.Games) <= 0 {
			return appError{"no games attached"}
		} else {
			gameName := ""
			if len(current.args) > 0 {
				gameName = current.args[0]
			}

			gamePath := getGameOrDefault(current.config.Games, gameName)
			manifest, err := loadGameManifest(gamePath)
			if err != nil {
				// TODO: embed error in appError
				return appError{fmt.Sprintf("game manifest not found")}
			}

			mods = manifest.Mods
		}

		//fmt.Printf("%-30s   %s\n", "NAME", "VERSION")
		if len(mods) > 0 {
			if useProfile {
				log("Mods installed in profile %s\n", profileName)
			}

			for name, version := range mods {
				log("%15s@%s\n", name, version)
			}
		} else {
			log("no mods installed\n")
		}

		return nil
	})
}

/* List Config Command */
type listConfigCommand struct{}

func (cmd *listConfigCommand) Usage() string {
	return ""
}

func (cmd *listConfigCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext, log logger) error {
		log("Settings: \n")
		log("  profile=%s\n", current.config.ActiveProfile)
		log("  channel=%s\n", current.config.ActiveRemote)
		return nil
	})
}

/* List Games Command */
type listGamesCommand struct{}

func (cmd *listGamesCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext, log logger) error {
		if len(current.config.Games) <= 0 {
			log("no games found.\n")
			return nil
		}

		log("%12s   %s\n", "ALIAS", "FOLDER")
		for alias, gameFolder := range current.config.Games {
			log("%12s   %s\n", alias, gameFolder)
		}

		return nil
	})
}

/* List Profiles Command */
type listProfilesCommand struct{}

func (cmd *listProfilesCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext, log logger) error {
		log("%15s   %s\n", "ID", "NAME")
		for _, p := range current.profiles {
			log("%15s   %s\n", p.Id, p.Name)
		}

		return nil
	})
}

/* List Channels Command */
type listChannelsCommand struct{}

func (cmd *listChannelsCommand) Usage() string {
	return ""
}

func (cmd *listChannelsCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext, log logger) error {
		format := "%15s  %10s   %s\n"
		log(format, "alias", "protocol", "endpoint")
		log(format, "==========", "========", "==========")
		for _, ch := range current.channels {
			log(format, ch.Alias, ch.Protocol, ch.Endpoint)
		}

		return nil
	})
}

func init() {
	list, _ := commandParser.AddCommand("list", "list manager data", "", &listCommand{})
	list.SubcommandsOptional = true
	list.AddCommand("mods", "lists the mods that are installed", "Lists the mods installed in the default or specified game.", &listModsCommand{})
	list.AddCommand("games", "lists the game folders attached", "Lists the games currently attached to the manager.", &listGamesCommand{})
	list.AddCommand("profiles", "lists available profiles", "List the available mod profiles", &listProfilesCommand{})
	list.AddCommand("config", "lists config settings", "Lists the current config settings.", &listConfigCommand{})
	list.AddCommand("channels", "list available channels", "Lists channels available in the manager.", &listChannelsCommand{})
}
