package sync

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("sync", Info)
}

var (
	Info = &cmd.Info{
		Use:   "sync",
		Short: "sync",
		Long:  "sync",
		SubCommands: []*cmd.Info{
			{
				Use:   "profile",
				Short: "profile",
				Long:  "profile",
				Run:   cmd.ExecutorFunc(syncProfile),
			},
			{
				Use:   "profiles",
				Short: "profiles",
				Long:  "profiles",
				Run:   cmd.ExecutorFunc(syncAllProfiles),
			},
		},
	}
)

func reportSyncResult(artifactName string, fromVersion string, toVersion string) {
	if fromVersion == toVersion {
		fmt.Printf("%s %-20s\n", artifactName, "OK")
	} else {
		fmt.Printf("%s %-20s->%s\n", artifactName, fromVersion, toVersion)
	}
}

func reportProfileSyncResult(p *Profile, from int32, to int32) {
	if from == to {
		fmt.Printf("%s @%d => no updates available\n", p.Name, from)
	} else {
		fmt.Printf("%s @%d => @%d\n", p.Name, from, to)
	}
}

/* Sync Profiles Command */
func syncAllProfiles(ctx context.Context, args []string) error {
	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
	}

	for _, p := range ctx.Profiles {
		if p.Source == nil {
			continue
		}

		from, to, err := p.Sync()
		if err != nil {
			log.Errorf("error syncing profile: %v", err)
			log.Println("couldn't sync with remote server.")
			return nil
		}

		err = p.Save()
		if err != nil {
			log.Errorf("error saving profile: %v", err)
			log.Println("couldn't save profile")
			return nil
		}

		reportProfileSyncResult(p, from, to)
	}

	return nil
}

/* Sync Profile Command */
func syncProfile(ctx context.Context, args []string) error {
	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
	}

	profile := ctx.Profile
	if len(args) > 0 {
		var ok bool
		profile, ok = ctx.Profile[args[0]]
		if !ok {
			log.Println("profile not found.")
			return nil
		}
	}

	if profile.Source == nil {
		log.Println("not a remote profile")
		return nil
	}

	from, to, err := profile.Sync()
	if err != nil {
		log.Errorf("error syncing profile: %v", err)
		log.Println("couldn't sync with remote server.")
		return nil
	}

	if err = profile.Save(); err != nil {
		log.Errorf("error saving profile: %v", err)
		log.Println("couldn't save profile")
		return nil
	}

	reportProfileSyncResult(profile, from, to)
	return nil
}
