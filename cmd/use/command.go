package use

import (
	"context"
	"errors"
	"log"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/configuration"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("use", Info)
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("profile name not specified")
	}

	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
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
		Use:   "use",
		Short: "use",
		Long:  "use",
		Run:   cmd.ExecutorFunc(run),
	}
)
