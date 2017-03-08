package profiles

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/cmd/cmdutils"
	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("profiles", Info)
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

			{
				Use:   "remove",
				Short: "remove a profile",
				Long:  "remove a profile",
				Run:   cmd.ExecutorFunc(removeProfile),
			},
			{
				Use:   "list",
				Short: "list profiles",
				Long:  "list profiles",
				Run:   cmd.ExecutorFunc(listProfiles),
			},
			{
				Use:   "export",
				Short: "export profile to file",
				Long:  "export a profile to a file",
				Run:   cmd.ExecutorFunc(exportProfile),
			},
		},
	}
)

/* Add Profiles Command */
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

/* Remove Profile Command */
func removeProfile(ctx context.Context, args []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("you must specify a profile")
	}

	profileId := args[0]
	if profile, err := m.RemoveProfile(profileId); err != nil {
		fmt.Printf("could not remove the profile: %v", err)
		return nil
	} else {
		fmt.Printf("%q has been removed\n", profile.Name)
	}

	return nil
}

/* List Profiles Command */
func listProfiles(ctx context.Context, _ []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%15s   %s\n", "ID", "NAME")
	for _, p := range m.Profiles() {
		fmt.Printf("%15s   %s\n", p.Id, p.Name)
	}

	return nil
}

/* Export Profile Command */
func exportProfile(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("missing profile name")
	} else if len(args) < 2 {
		return errors.New("missing export path")
	}

	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	profileId := args[0]
	profile, ok := m.Profiles()[profileId]
	if !ok {
		return fmt.Errorf("[%s] not found\n", profileId)
	}

	profilePath := args[1]
	if err := manager.SaveProfile(m.Fs(), profilePath, profile); err != nil {
		fmt.Printf("error saving profile: %v", err)
		return nil
	}

	fmt.Printf("[%s] exported to: %s\n", profile.Id, profilePath)
	return nil
}
