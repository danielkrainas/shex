package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/cmd/cmdutils"
)

func init() {
	cmd.Register("cache", Info)
}

var (
	Info = &cmd.Info{
		Use:   "cache",
		Short: "cache",
		Long:  "cache",
		SubCommands: []*cmd.Info{
			&cmd.Info{
				Use:   "clean",
				Short: "clean",
				Long:  "clean",
				Run:   cmd.ExecutorFunc(cleanCache),
			},
		},
	}
)

func cleanCache(ctx context.Context, _ []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	if err := m.ClearCache(); err != nil {
		fmt.Printf("error clearing cache: %v", err)
		return nil
	}

	log.Printf("cache cleared")
	return nil
}