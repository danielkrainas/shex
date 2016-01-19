package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type addCommand struct{}

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
		err := profile.saveTo(profilePath)
		if err != nil {
			return appError{err, fmt.Sprintf("Could not save to: %s", profilePath)}
		}

		log.Printf("[%s] created at: %s\n", profile.Id, profilePath)
		return nil
	})
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

/* Add Channel Command */
type addChannelCommand struct {
	Protocol string `short:"p" long:"protocol" description:"set the protocol to use with the channel"`
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
	add, _ := commandParser.AddCommand("add", "", "", &addCommand{})
	add.AddCommand("profile", "", "", &addProfileCommand{})
	add.AddCommand("game", "", "", &addGameCommand{})
	add.AddCommand("channel", "", "", &addChannelCommand{})
}
