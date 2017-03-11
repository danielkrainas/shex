package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	_ "github.com/danielkrainas/shex/cmd/api"
	"github.com/danielkrainas/shex/cmd/registry"
	_ "github.com/danielkrainas/shex/cmd/version"
)

var appVersion string

const DEFAULT_VERSION = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = DEFAULT_VERSION
	}

	rand.Seed(time.Now().Unix())
	ctx := acontext.WithVersion(acontext.Background(), appVersion)

	dispatch := cmd.CreateDispatcher(ctx, registry.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
