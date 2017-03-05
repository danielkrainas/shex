package profiles

import (
	"context"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/cmd/cmdutils"
	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("profiles", Info)
}

func addProfile(ctx context.Context, args []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	profileId := args[0]
	profilePath := ""
	if len(args) > 1 {
		profilePath = args[1]
	}

	var profile *v1.Profile
	if profilePath != "" {
		if p, err := manager.LoadProfile(m.Fs(), profilePath); err != nil {
			return err
		} else {
			profile = p
		}
	} else {
		profile = v1.NewProfile(profileId)
	}

	if err := m.AddProfile(profile); err != nil {
		fmt.Printf("error saving profile: %v", err)
		fmt.Printf("Could not save to: %s", profilePath)
		return nil
	}

	fmt.Printf("[%s] created at: %s\n", profile.Id, profilePath)
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "profiles",
		Short: "",
		Long:  "",
		SubCommands: []*cmd.Info{
			{
				Use:   "add",
				Short: "add a profile",
				Long:  "add a profile",
				Run:   cmd.ExecutorFunc(addProfile),
			},
		},
	}
)
