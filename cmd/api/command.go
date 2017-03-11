package api

import (
	"context"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/api/server"
	"github.com/danielkrainas/shex/registry/configuration"
)

func init() {
	cmd.Register("api", Info)
}

func run(ctx context.Context, args []string) error {
	config, err := configuration.Resolve(args)
	if err != nil {
		return err
	}

	s, err := server.New(ctx, config)
	if err != nil {
		return err
	}

	return s.ListenAndServe()
}

var (
	Info = &cmd.Info{
		Use:   "api",
		Short: "run the api server",
		Long:  "Run the api server.",
		Run:   cmd.ExecutorFunc(run),
	}
)
