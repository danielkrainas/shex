package uninstall

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
	cmd.Register("uninstall", Info)
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("profile name not specified")
	}

	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
	}

	isSelf := args[0] == "self"
	if isSelf {
		installPath := ""
		if len(args) > 1 {
			installPath = args[1]
		}

		if err := uninstallSelf(installPath); err != nil {
			log.Errorf("error uninstalling self: %v", err)
			log.Println("Could not uninstall. Depending on your system's configuration, you may need to run the uninstall again as an administrator.")
		}

		return nil
	}

	name := args[0]
	gamePath := game.GetGameOrDefault(ctx.Config.Games, name)
	mod, err := game.UninstallMod(ctx.Config, gamePath, ctx.Profile, name)
	if err != nil {
		log.Errorf("error uninstalling mod: %v", err)
		log.Println("Could not uninstall mod")
		return nil
	}

	log.Printf("%s@%s uninstalled", mod.Name, mod.SemVersion)
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "uninstall",
		Short: "u",
		Long:  "uninstall",
		Run:   cmd.ExecutorFunc(run),
	}
)
