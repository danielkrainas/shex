package add

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("add", Info)
}

func addWrapper(fn func(*manager.ExecutionContext, []string) error) func(context.Context, []string) error {
	return func(ctx context.Context, args []string) error {
		ctx, err := manager.Context(ctx, "")
		if err != nil {
			return err
		}

		return fn(ctx, args)
	}
}

var (
	Info = &cmd.Info{
		Use:   "add",
		Short: "add",
		Long:  "add",
		Run:   cmd.ExecutorFunc(run),
		Commands: []*cmd.Info{
			{
				Use:   "profile",
				Short: "profile",
				Long:  "profile",
				Run:   cmd.ExecutorFunc(addWrapper(addProfile)),
			},
			{
				Use:   "game",
				Short: "game",
				Long:  "game",
				Run:   cmd.ExecutorFunc(addWrapper(addGame)),
			},
			{
				Use:   "channel",
				Short: "channel",
				Long:  "channel",
				Run:   cmd.ExecutorFunc(addWrapper(addChannel)),
				Flags: []*cmd.Flag{
					{
						Long:        "protocol",
						Short:       "p",
						Description: "set the protocol to use with the channel",
					},
				},
			},
		},
	}
)

func addProfile(ctx *manager.ExecutionContext, args []string) error {
	profileId := args[0]
	profilePath := ""
	if len(args) > 1 {
		profilePath = args[1]
	} else {
		profilePath = path.Join(ctx.Config.ProfilesPath, profileId+".json")
	}

	if profile, ok := ctx.Profiles[profileId]; ok {
		log.Printf("[%s] already exists\n", profile.Id)
		return nil
	}

	profile := manager.CreateProfile(profileId)
	if err := profile.SaveTo(profilePath); err != nil {
		log.Errorf("error saving profile: %v", err)
		log.Printf("Could not save to: %s", profilePath)
		return nil
	}

	log.Printf("[%s] created at: %s\n", profile.Id, profilePath)
	return nil
}

func addGame(ctx *manager.ExecutionContext, args []string) error {
	if len(args) < 1 {
		return errors.New("invalid game path")
	}

	alias := args[0]
	var gamePath string
	var err error
	if len(args) < 2 {
		gamePath = alias
		alias = manager.DefaultGameName
		log.Println("No alias specified, assuming \"default\"")
	} else {
		gamePath, err = filepath.Abs(args[1])
		if err != nil {
			return errors.New("couldn't resolve path: " + args[1])
		}
	}

	if ctx.Config.Games.Exists(alias) {
		log.Errorf("the alias %q is already in use", alias)
		return nil
	}

	if err = game.AttachGameFolder(ctx.Config, alias, gamePath); err != nil {
		log.Errorf("error attaching game folder: %v", err)
		log.Printf("could not attach game %q at: %s", alias, gamePath)
		return nil
	}

	if err = configuration.Save(ctx.Config, ctx.HomePath); err != nil {
		log.Errorf("error saving manager config: %v", err)
		log.Printf("could not save config: %s", ctx.HomePath)
		return nil
	}

	gamePath = ctx.Config.Games[alias]
	log.Printf("added %s as \"%s\"\n", gamePath, alias)
	return nil
}

func addChannel(ctx *manager.ExecutionContext, args []string) error {
	if len(args) < 2 {
		return errors.New("alias and/or endpoint missing")
	}

	alias := strings.ToLower(args[0])
	endpoint := args[1]
	if c, ok := ctx.Config.Channels[alias]; ok {
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

	channelPath := filepath.Join(ctx.Config.ChannelsPath, c.Alias+".json")
	if err := c.SaveTo(channelPath); err != nil {
		log.Errorf("error saving channel: %v", err)
		log.Println("couldn't save channel: ", channelPath)
		return nil
	}

	log.Printf("channel added: %s => %s\n", c.Alias, channelPath)
	return nil
}
