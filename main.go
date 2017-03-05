package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/danielkrainas/gobag/cmd"
	"github.com/danielkrainas/gobag/context"

	_ "github.com/danielkrainas/shex/cmd/games"
	_ "github.com/danielkrainas/shex/cmd/profiles"
	//_ "github.com/danielkrainas/shex/cmd/add"
	//_ "github.com/danielkrainas/shex/cmd/clean"
	//_ "github.com/danielkrainas/shex/cmd/export"
	//_ "github.com/danielkrainas/shex/cmd/install"
	//_ "github.com/danielkrainas/shex/cmd/list"
	//_ "github.com/danielkrainas/shex/cmd/remove"
	"github.com/danielkrainas/shex/cmd/root"
	//_ "github.com/danielkrainas/shex/cmd/set"
	//_ "github.com/danielkrainas/shex/cmd/sync"
	//_ "github.com/danielkrainas/shex/cmd/uninstall"
	//_ "github.com/danielkrainas/shex/cmd/use"
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

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
