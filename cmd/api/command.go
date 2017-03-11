package api

import (
	"context"
	"fmt"

	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"
)

func init() {
	cmd.Register("api", Info)
}

func run(ctx context.Context, args []string) error {

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "api",
		Short: "run the api server",
		Long:  "run the api server",
		Run:   cmd.ExecutorFunc(run),
	}
)
