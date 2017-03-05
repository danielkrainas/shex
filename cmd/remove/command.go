package remove

import (
	"context"
	"errors"
	"log"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/manager"
)

func init() {
	cmd.Register("remove", Info)
}

var (
	Info = &cmd.Info{
		Use:   "remove",
		Short: "remove",
		Long:  "remove",
		SubCommands: []*cmd.Info{
			{
				Use:   "channel",
				Short: "channel",
				Long:  "channel",
				Run:   cmd.ExecutorFunc(removeChannel),
			},
		},
	}
)

/* Remove Channel Command */
func removeChannel(parent context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("you must specify a channel alias")
	}

	ctx, err := manager.Context(parent, "")
	if err != nil {
		return err
	}

	alias := args[0]
	var channel *manager.Channel
	ok := false
	if alias == "default" && ctx.Config.IncludeDefaultChannel {
		channel = manager.DefaultChannel
		ok = true
	} else {
		channel, ok = ctx.Channels[alias]
	}

	if !ok {
		log.Println("channel not found")
		return nil
	}

	if channel == manager.DefaultChannel {
		ctx.Config.IncludeDefaultChannel = false
		if err := manager.SaveConfig(ctx.Config, ctx.HomePath); err != nil {
			log.Printf("error saving configuration: %v", err)
			log.Println("couldn't save configuration")
			return nil
		}
	} else if err := channel.Remove(); err != nil {
		log.Printf("error removing channel: %v", err)
		log.Println("couldn't remove channel")
		return nil
	}

	log.Printf("channel removed: %s => %s\n", channel.Alias, channel.Endpoint)
	return nil
}
