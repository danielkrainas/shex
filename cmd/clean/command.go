package clean

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
	cmd.Register("clean", Info)
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("resource type not specified")
	}

	ctx, err := manager.Context(ctx, "")
	if err != nil {
		return err
	}

	targetType := args[0]
	if targetType != "cache" {
		return errors.New("resource must be one of: `cache`")
	}

	targetPath := filepath.Join(current.homePath, ctx.Config.CachePath)
	if err := fsutils.ClearDirectory(targetPath); err != nil {
		return fmt.Errorf("error clearing %q: %v", targetPath, err)
	}

	log.Printf("cleared %s => %s", targetType, targetPath)
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "clean",
		Short: "clean",
		Long:  "clean",
		Run:   cmd.ExecutorFunc(run),
	}
)
