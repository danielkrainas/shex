package list

import (
	"context"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/manager"
	"github.com/danielkrainas/shex/mods"
)

func init() {
	cmd.Register("list", Info)
}

func listWrapper(fn func(*manager.ExecutionContext, []string) error) func(context.Context, []string) error {
	return func(parent context.Context, args []string) error {
		ctx, err := manager.Context(parent, "")
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
		SubCommands: []*cmd.Info{
			{
				Use:   "mods",
				Short: "mods",
				Long:  "mods",
				Run:   cmd.ExecutorFunc(listWrapper(listMods)),
			},
			{
				Use:   "config",
				Short: "config",
				Long:  "config",
				Run:   cmd.ExecutorFunc(listWrapper(listConfig)),
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
	var list v1.ModList
	if useProfile {
		profileName = ctx.Value("flags.profile").(string)
		if len(profileName) > 0 {
			selectedProfile, ok := ctx.Profiles[profileName]
			if !ok {
				return fmt.Errorf("profile not found: %q", profileName)
			}

			list = selectedProfile.Mods
		} else {
			profileName = ctx.Profile().Name
			list = ctx.Profile().Mods
		}
	} else if len(ctx.Config.Games) <= 0 {
		log.Println("no games attached")
		return nil
	} else {
		gameName := ""
		if len(args) > 0 {
			gameName = args[0]
		}

		gamePath := manager.GetGameOrDefault(ctx.Config.Games, gameName)
		manifest, err := mods.LoadGameManifest(gamePath)
		if err != nil {
			log.Printf("error loading game manifest: %v", err)
			log.Println("game manifest not found or invalid")
			return nil
		}

		list = manifest.Mods
	}

	//fmt.Printf("%-30s   %s\n", "NAME", "VERSION")
	if len(list) > 0 {
		if useProfile {
			log.Printf("Mods installed in profile %s\n", profileName)
		}

		for name, version := range list {
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
