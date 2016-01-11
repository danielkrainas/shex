package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type addCommand struct{}

func (cmd *addCommand) Usage() string {
	return ""
}

type addProfileCommand struct{}

func (cmd *addProfileCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext) error {
		profileId := args[0]
		profilePath := ""
		if len(args) > 1 {
			profilePath = args[1]
		} else {
			profilePath = path.Join(current.config.ProfilesPath, profileId+".json")
		}

		if profile, ok := current.profiles[profileId]; ok {
			return appError{nil, fmt.Sprintf("[%s] already exists\n", profile.Id)}
		}

		profile := createProfile(profileId)
		err := saveProfile(profile, profilePath)
		if err != nil {
			return appError{err, fmt.Sprintf("Could not save to: %s", profilePath)}
		}

		log.Printf("[%s] created at: %s\n", profile.Id, profilePath)
		return nil
	})
}

func (cmd *addProfileCommand) Usage() string {
	return "<id> [path]"
}

/* Add Game Command */
type addGameCommand struct{}

func (cmd *addGameCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext) error {
		var alias string
		var gamePath string
		var err error

		alias = args[0]
		if len(args) < 2 {
			gamePath = alias
			alias = DefaultGameName
			log.Printf("No alias specified, assuming \"default\"")
		} else {
			gamePath, err = filepath.Abs(args[1])
			if err != nil {
				return appError{err, "couldn't resolve path: " + args[1]}
			}
		}

		if _, ok := current.config.Games[alias]; ok {
			return appError{nil, fmt.Sprintf("The alias \"%s\" is already in use\n", alias)}
		}

		if err = attachGameFolder(current.config, alias, gamePath); err != nil {
			return appError{err, fmt.Sprintf("Could not attach game \"%s\" at: %s", alias, gamePath)}
		}

		if err = saveManagerConfig(current.config, current.homePath); err != nil {
			return appError{err, fmt.Sprintf("Could not save config: %s", current.homePath)}
		}

		gamePath = current.config.Games[alias]
		log.Printf("added %s as \"%s\"\n", gamePath, alias)
		return nil
	})
}

func (cmd *addGameCommand) Usage() string {
	return "<alias> <game_path>"
}

/* Add Channel Command */
type addChannelCommand struct {
	Protocol string `short:"p" long:"protocol" description:"set the protocol to use with the channel"`
}

func (cmd *addChannelCommand) Usage() string {
	return "<alias> <endpoint> [options]"
}

func (cmd *addChannelCommand) Execute(args []string) error {
	if len(args) < 2 {
		return usageError{}
	}

	return runInContext(func(current *executionContext) error {
		alias := strings.ToLower(args[0])
		endpoint := args[1]
		if c, ok := current.config.channels[alias]; ok {
			log.Printf("Overriding %s (=%s:%s)\n", alias, c.Protocol, c.Endpoint)
		}

		c := &Channel{
			Alias:    alias,
			Endpoint: endpoint,
		}

		if len(cmd.Protocol) < 1 {
			cmd.Protocol = "http"
		}

		if cmd.Protocol != "http" && cmd.Protocol != "https" {
			return appError{nil, "unknown protocol: " + cmd.Protocol}
		} else {
			c.Protocol = cmd.Protocol
		}

		channelPath := filepath.Join(current.config.ChannelsPath, c.Alias+".json")
		if err := c.saveTo(channelPath); err != nil {
			return appError{err, "couldn't save channel: " + channelPath}
		}

		log.Printf("channel added: %s => %s\n", c.Alias, channelPath)
		return nil
	})
}

func init() {
	add, _ := commandParser.AddCommand("add", "adds ", "", &addCommand{})
	add.AddCommand("profile", "add a profile config", "Creates a new mod profile with the specified id. If a path argument is supplied, the profile won't be imported and will be saved to the path specified.", &addProfileCommand{})
	add.AddCommand("game", "add a game folder", "Adds the game folder at the specified location to the manager. <alias> may be omitted and \"default\" will be assumed.", &addGameCommand{})
	add.AddCommand("channel", "add a remote channel", "Adds a remote channel to the manager.", &addChannelCommand{})
}
