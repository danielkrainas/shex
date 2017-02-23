package remove

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"
)

func init() {
	cmd.Register("remove", Info)
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("profile name not specified")
	}

	newProfileName := args[0]
	if newProfileName != ctx.Config.ActiveProfile {
		newProfile := ctx.Profiles[newProfileName]
		ctx.Config.ActiveProfile = newProfile.Id
		if err := configuration.Save(ctx.Config, ctx.HomePath); err != nil {
			return err
		}

		log.Printf("active profile set to: %s\n", newProfile.Name)
	} else {
		log.Printf("profile already active")
	}

	return nil
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
func removeProfile(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("you must specify a profile")
	}

	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
	}

	profileId := args[0]
	profile, ok := current.profiles[profileId]
	if !ok {
		return appError{nil, fmt.Sprintf("Could not find the profile %q", profileId)}
	}

	if err := profile.drop(); err != nil {
		return err
	}

	log.Printf("%q has been removed\n", profile.Name)
	return nil
}

/* Remove Game Command */
func removeGame(ctx context.Context, args []string) error {
	if len(args) < 0 {
		return errors.New("you must specify a game alias")
	}

	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
	}

	alias := args[0]
	gamePath, ok := current.config.Games[alias]
	if !ok {
		log.Printf("game %q does not exist.", alias)
		return nil
	}

	err := game.DetachGameFolder(current.config, alias)
	if err != nil {
		log.Errorf("error removing game folder: %v", err)
		log.Println("Could not remove game from manager")
		return nil
	}

	err = configuration.Save(ctx.Config, ctx.HomePath)
	if err != nil {
		log.Errorf("error saving configuration: %v", err)
		log.Println("could not save configuration")
		return nil
	}

	log.Printf("game removed: %s => %s", alias, gamePath)
	return nil
}

/* Remove Channel Command */
func removeChannel(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("you must specify a channel alias")
	}

	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
	}

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
		log.Println("channel not found")
		return nil
	}

	if channel == defaultChannel {
		current.config.IncludeDefaultChannel = false
		if err := saveManagerConfig(current.config, current.homePath); err != nil {
			log.Errorf("error saving configuration: %v", err)
			log.Println("couldn't save configuration")
			return nil
		}
	} else if err := channel.remove(); err != nil {
		log.Errorf("error removing channel: %v", err)
		log.Println("couldn't remove channel")
		return nil
	}

	log.Printf("channel removed: %s => %s\n", channel.Alias, channel.Endpoint)
	return nil
}
