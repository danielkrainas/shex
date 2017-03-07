package channels

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/cmd/cmdutils"
	"github.com/danielkrainas/shex/manager"
	"github.com/danielkrainas/shex/mods"
)

func init() {
	cmd.Register("channels", Info)
}

var (
	Info = &cmd.Info{
		Use:   "channels",
		Short: "",
		Long:  "",
		SubCommands: []*cmd.Info{
			{
				Use:   "add",
				Short: "add a channel",
				Long:  "add a channel",
				Run:   cmd.ExecutorFunc(addChannel),
				Flags: []*cmd.Flag{
					{
						Long:        "protocol",
						Short:       "p",
						Description: "set the protocol to use with the channel",
					},
				},
			},

			{
				Use:   "remove",
				Short: "remove a channel",
				Long:  "remove a channel",
				Run:   cmd.ExecutorFunc(removeChannel),
			},
			{
				Use:   "list",
				Short: "list channels",
				Long:  "list channels",
				Run:   cmd.ExecutorFunc(listChannels),
			},
		},
	}
)

/* Add Channel Command */
func addChannel(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.New("alias and/or endpoint missing")
	}

	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	alias := strings.ToLower(args[0])
	endpoint := args[1]
	if ch, ok := m.Channels()[alias]; ok {
		fmt.Printf("Overriding %s (=%s:%s)\n", alias, ch.Protocol, ch.Endpoint)
	}

	ch := &mods.Channel{
		Alias:    alias,
		Endpoint: endpoint,
	}

	protocol := ctx.Value("flags.protocol").(string)
	if len(protocol) < 1 {
		protocol = "http"
	}

	if protocol != "http" && protocol != "https" {
		fmt.Printf("unknown protocol: %s", protocol)
		return nil
	} else {
		ch.Protocol = protocol
	}

	if err := m.AddChannel(ch); err != nil {
		fmt.Printf("error adding channel: %v", err)
		return nil
	}

	fmt.Printf("channel added: %s", ch.Alias)
	return nil
}

/* Remove Channel Command */
func removeChannel(ctx context.Context, args []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("you must specify a channel alias")
	}

	alias := args[0]
	var ch *mods.Channel
	if ch, err = m.RemoveChannel(alias); err != nil {
		fmt.Printf("error removing channel: %v", err)
		return nil
	}

	if ch == manager.DefaultChannel {
		if err := m.SaveConfig(); err != nil {
			fmt.Printf("error saving config: %v", err)
			return nil
		}
	}

	fmt.Printf("channel removed: %s => %s\n", ch.Alias, ch.Endpoint)
	return nil
}

/* List Channels Command */
func listChannels(ctx context.Context, _ []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	format := "%15s  %10s   %s\n"
	fmt.Printf(format, "alias", "protocol", "endpoint")
	fmt.Printf(format, "==========", "========", "==========")
	for _, ch := range m.Channels() {
		fmt.Printf(format, ch.Alias, ch.Protocol, ch.Endpoint)
	}

	return nil
}