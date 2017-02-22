package install

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/fsutils"
	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("install", Info)
}

func run(ctx *manager.ExecutionContext, args []string) error {
	if len(args) < 1 {
		return errors.New("must specify a target")
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

		if err := self.InstallSelf(installPath); err != nil {
			log.Errorf("error installing self: %v", err)
			log.Printf("could not install locally. Depending on your system's configuration, you may need to run the install as an administrator.")
		}

		return nil
	}

	token := game.ParseNameVersionToken(args[0])
	gamePath := ctx.Config.Games.GameOrDefault("")
	mod, err := game.InstallMod(ctx.Config, gamePath, ctx.Profile, token)
	if err != nil {
		log.Printf("error installing mod: %v", err)
		log.Printf("could not install mod: %v", err)
		return nil
	}

	log.Printf("%s@%s installed at %s\n", mod.Name, mod.SemVersion, mod.source)
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "install",
		Short: "i",
		Long:  "install",
		Run:   cmd.ExecutorFunc(run),
	}
)
