package use

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("export", Info)
}

func run(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("profile name not specifiedg")
	}

	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	profileId := args[0]
	profile, ok := ctx.Profiles[profileId]
	if !ok {
		return fmt.Errorf("[%s] not found\n", profileId)
	}

	profilePath := args[1]
	if err := manager.SaveProfileTo(profile, profilePath); err != nil {
		log.Printf("error saving profile: %v", err)
		return nil
	}

	log.Printf("[%s] exported to: %s\n", profile.Id, profilePath)
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "export",
		Short: "export",
		Long:  "export",
		SubCommands: []*cmd.Info{
			{
				Use:   "profile",
				Short: "profile",
				Long:  "profile",
				Run:   cmd.ExecutorFunc(run),
			},
		},
	}
)
