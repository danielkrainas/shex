package remove

import (
	"context"
	"errors"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("remove", Info)
}

var (
	Info = &cmd.Info{
		Use:   "remove",
		Short: "remove",
		Long:  "remove",
		SubCommands: []*cmd.Info{
			{
				Use:   "game",
				Short: "game",
				Long:  "game",
				Run:   cmd.ExecutorFunc(removeGame),
			},
			{
				Use:   "profile",
				Short: "profile",
				Long:  "profile",
				Run:   cmd.ExecutorFunc(removeProfile),
			},
			{
				Use:   "channel",
				Short: "channel",
				Long:  "channel",
				Run:   cmd.ExecutorFunc(removeChannel),
			},
		},
	}
)

/* Remove Profile Command */
func removeProfile(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("you must specify a profile")
	}

	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	profileId := args[0]
	profile, ok := ctx.Profiles[profileId]
	if !ok {
		log.Printf("Could not find the profile %q", profileId)
		return nil
	}

	if err := manager.DropProfile(profile); err != nil {
		return err
	}

	log.Printf("%q has been removed\n", profile.Name)
	return nil
}

/* Remove Game Command */
func removeGame(parent context.Context, args []string) error {
	if len(args) < 0 {
		return errors.New("you must specify a game alias")
	}

	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	alias := args[0]
	gamePath, ok := ctx.Config.Games[alias]
	if !ok {
		log.Printf("game %q does not exist.", alias)
		return nil
	}

	ctx.Config.Games.Detach(alias)
	if err := manager.SaveConfig(ctx.Config, ctx.HomePath); err != nil {
		log.Printf("error saving configuration: %v", err)
		log.Println("could not save configuration")
		return nil
	}

	log.Printf("game removed: %s => %s", alias, gamePath)
	return nil
}

/* Remove Channel Command */
func removeChannel(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("you must specify a channel alias")
	}

	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	alias := args[0]
	var channel *manager.Channel
	ok := false
	if alias == "default" && ctx.Config.IncludeDefaultChannel {
		channel = manager.DefaultChannel
		ok = true
	} else {
		channel, ok = ctx.Channels[alias]
	}

	if !ok {
		log.Println("channel not found")
		return nil
	}

	if channel == manager.DefaultChannel {
		ctx.Config.IncludeDefaultChannel = false
		if err := manager.SaveConfig(ctx.Config, ctx.HomePath); err != nil {
			log.Printf("error saving configuration: %v", err)
			log.Println("couldn't save configuration")
			return nil
		}
	} else if err := channel.Remove(); err != nil {
		log.Printf("error removing channel: %v", err)
		log.Println("couldn't remove channel")
		return nil
	}

	log.Printf("channel removed: %s => %s\n", channel.Alias, channel.Endpoint)
	return nil
}
