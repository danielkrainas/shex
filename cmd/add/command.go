package add

import (
	"context"
	"errors"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("add", Info)
}

func addWrapper(fn func(*manager.ExecutionContext, []string) error) func(context.Context, []string) error {
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
		Use:   "add",
		Short: "add",
		Long:  "add",
		SubCommands: []*cmd.Info{
			{
				Use:   "game",
				Short: "game",
				Long:  "game",
				Run:   cmd.ExecutorFunc(addWrapper(addGame)),
			},
			{
				Use:   "channel",
				Short: "channel",
				Long:  "channel",
				Run:   cmd.ExecutorFunc(addWrapper(addChannel)),
				Flags: []*cmd.Flag{
					{
						Long:        "protocol",
						Short:       "p",
						Description: "set the protocol to use with the channel",
					},
				},
			},
		},
	}
)

func addGame(ctx *manager.ExecutionContext, args []string) error {
	if len(args) < 1 {
		return errors.New("invalid game path")
	}

	alias := args[0]
	var gamePath string
	var err error
	if len(args) < 2 {
		gamePath = alias
		alias = manager.DefaultGameName
		log.Println("No alias specified, assuming \"default\"")
	} else {
		gamePath, err = filepath.Abs(args[1])
		if err != nil {
			return errors.New("couldn't resolve path: " + args[1])
		}
	}

	if ctx.Config.Games.HasAlias(alias) {
		log.Printf("the alias %q is already in use", alias)
		return nil
	}

	ctx.Config.Games.Attach(alias, gamePath)
	if err = manager.SaveConfig(ctx.Config, ctx.HomePath); err != nil {
		log.Printf("error saving manager config: %v", err)
		log.Printf("could not save config: %s", ctx.HomePath)
		return nil
	}

	gamePath = ctx.Config.Games[alias]
	log.Printf("added %s as \"%s\"\n", gamePath, alias)
	return nil
}

func addChannel(ctx *manager.ExecutionContext, args []string) error {
	if len(args) < 2 {
		return errors.New("alias and/or endpoint missing")
	}

	alias := strings.ToLower(args[0])
	endpoint := args[1]
	if c, ok := ctx.Channels[alias]; ok {
		log.Printf("Overriding %s (=%s:%s)\n", alias, c.Protocol, c.Endpoint)
	}

	c := &manager.Channel{
		Alias:    alias,
		Endpoint: endpoint,
	}

	protocol := ctx.Value("flags.protocol").(string)
	if len(protocol) < 1 {
		protocol = "http"
	}

	if protocol != "http" && protocol != "https" {
		log.Printf("unknown protocol: %s", protocol)
		return nil
	} else {
		c.Protocol = protocol
	}

	channelPath := filepath.Join(ctx.Config.ChannelsPath, c.Alias+".json")
	if err := c.SaveTo(channelPath); err != nil {
		log.Printf("error saving channel: %v", err)
		log.Println("couldn't save channel: ", channelPath)
		return nil
	}

	log.Printf("channel added: %s => %s\n", c.Alias, channelPath)
	return nil
}
