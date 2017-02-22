package list

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
	cmd.Register("list", Info)
}

func listWrapper(fn func(*manager.ExecutionContext, []string) error) func(context.Context, []string) error {
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
		Use:   "list",
		Short: "list",
		Long:  "list",
		Run:   cmd.ExecutorFunc(run),
		Commands: []*cmd.Info{
			{
				Use:   "mods",
				Short: "mods",
				Long:  "mods",
				Run:   cmd.ExecutorFunc(listWrapper(listMods)),
			},

			{
				Use:   "games",
				Short: "games",
				Long:  "games",
				Run:   cmd.ExecutorFunc(listWrapper(listGames)),
			},
			{
				Use:   "profiles",
				Short: "profiles",
				Long:  "profiles",
				Run:   cmd.ExecutorFunc(listWrapper(listProfiles)),
			},
			{
				Use:   "config",
				Short: "config",
				Long:  "config",
				Run:   cmd.ExecutorFunc(listWrapper(listConfig)),
			},
			{
				Use:   "channels",
				Short: "channels",
				Long:  "channels",
				Run:   cmd.ExecutorFunc(listWrapper(listChannels)),
			},
		},
	}
)

/* List Mods Command
type listModsCommand struct {
	Profile string `short:"p" long:"profile" description:"display mods installed in a profile"`
} */

func listMods(ctx *manager.ExecutionContext, args []string) error {
	profileName := acontext.GetStringValue(ctx, "flags.profile")
	useProfile := profileName != ""
	var mods ModList
	if useProfile {
		profileName = cmd.Profile
		if len(cmd.Profile) > 0 {
			selectedProfile, ok := ctx.Profiles[profileName]
			if !ok {
				return fmt.Errorf("profile not found: %q", profileName)
			}

			mods = selectedProfile.Mods
		} else {
			profileName = ctx.Profile.Name
			mods = ctx.Profile.Mods
		}
	} else if len(ctx.Config.Games) <= 0 {
		log.Println("no games attached")
		return nil
	} else {
		gameName := ""
		if len(args) > 0 {
			gameName = args[0]
		}

		gamePath := getGameOrDefault(ctx.Config.Games, gameName)
		manifest, err := loadGameManifest(gamePath)
		if err != nil {
			log.Errorf("error loading game manifest: %v", err)
			log.Println("game manifest not found or invalid")
			return nil
		}

		mods = manifest.Mods
	}

	//fmt.Printf("%-30s   %s\n", "NAME", "VERSION")
	if len(mods) > 0 {
		if useProfile {
			log.Printf("Mods installed in profile %s\n", profileName)
		}

		for name, version := range mods {
			log.Printf("%15s@%s\n", name, version)
		}
	} else {
		log.Printf("no mods installed\n")
	}

	return nil
}

/* List Config Command */
func listConfig(ctx *manager.ExecutionContext, _ []string) error {
	log.Printf("Settings: \n")
	log.Printf("    profile=%s\n", ctx.Config.ActiveProfile)
	log.Printf("    channel=%s\n", ctx.Config.ActiveRemote)
	return nil
}

/* List Games Command */
func listGames(ctx *manager.ExecutionContext, _ []string) error {
	if len(ctx.Config.Games) <= 0 {
		log.Printf("no games found.\n")
		return nil
	}

	log.Printf("%12s   %s\n", "ALIAS", "FOLDER")
	for alias, gameFolder := range ctx.Config.Games {
		log.Printf("%12s   %s\n", alias, gameFolder)
	}

	return nil
}

/* List Profiles Command */
func listProfiles(ctx *manager.ExecutionContext, _ []string) error {
	log.Printf("%15s   %s\n", "ID", "NAME")
	for _, p := range ctx.Profiles {
		log.Printf("%15s   %s\n", p.Id, p.Name)
	}

	return nil
}

/* List Channels Command */
func listChannels(ctx *manager.ExecutionContext, _ []string) error {
	format := "%15s  %10s   %s\n"
	log.Printf(format, "alias", "protocol", "endpoint")
	log.Printf(format, "==========", "========", "==========")
	for _, ch := range ctx.Channels {
		log.Printf(format, ch.Alias, ch.Protocol, ch.Endpoint)
	}

	return nil
}
