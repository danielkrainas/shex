package set

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("set", Info)
}

func run(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("resource type not specified")
	}

	if len(args) < 1 {
		return errors.New("you must specify a setting key")
	} else if len(args) < 2 {
		return errors.New("setting value missing")
	}

	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	target := args[0]
	value := args[1]
	switch target {
	case "profile":
		if _, ok := ctx.Profiles[value]; ok {
			ctx.Config.ActiveProfile = value
		} else {
			log.Println("profile not found")
			return nil
		}

		break

	case "channel":
		if _, ok := ctx.Channels[value]; ok {
			ctx.Config.ActiveRemote = value
		} else {
			log.Println("channel not found")
			return nil
		}

		break

		/*case "game":
		if game, ok := current.config.Games[value]; ok {
			current.config.ActiveGame = value
		} else {
			return appError{"game not found"}
		}

		break*/
	default:
		return fmt.Errorf("unknown setting key: %s", target)
	}

	if err := manager.SaveConfig(ctx.Config, ctx.HomePath); err != nil {
		log.Printf("error saving config: %v", err)
		log.Println("couldn't save configuration")
		return nil
	}

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "set",
		Short: "set",
		Long:  "set",
		Run:   cmd.ExecutorFunc(run),
	}
)
