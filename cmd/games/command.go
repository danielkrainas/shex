package games

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/shex/cmd/cmdutils"
	"github.com/danielkrainas/shex/manager"
	"github.com/danielkrainas/shex/mods"
)

func init() {
	cmd.Register("games", Info)
}

var (
	Info = &cmd.Info{
		Use:   "games",
		Short: "",
		Long:  "",
		SubCommands: []*cmd.Info{
			{
				Use:   "add",
				Short: "add a game folder",
				Long:  "add a game folder",
				Run:   cmd.ExecutorFunc(addGame),
			},

			{
				Use:   "remove",
				Short: "remove a game folder",
				Long:  "remove a game folder",
				Run:   cmd.ExecutorFunc(removeGame),
			},
			{
				Use:   "list",
				Short: "list",
				Long:  "list",
				Run:   cmd.ExecutorFunc(listGames),
			},
		},
	}
)

/* Add Game Command */
func addGame(ctx context.Context, args []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("invalid game path")
	}

	alias := args[0]
	var gamePath string
	if len(args) < 2 {
		gamePath = alias
		alias = manager.DefaultGameName
		fmt.Println(`No alias specified, assuming "default"`)
	} else {
		gamePath, err = filepath.Abs(args[1])
		if err != nil {
			return errors.New("couldn't resolve path: " + args[1])
		}
	}

	if err := m.AddGame(alias, mods.GameDir(gamePath)); err != nil {
		fmt.Printf("error adding game: %v", err)
		return nil
	}

	if err := m.SaveConfig(); err != nil {
		fmt.Printf("error saving config: %v", err)
		return nil
	}

	fmt.Printf("added %s as %q", gamePath, alias)
	return nil
}

/* Remove Game Command */
func removeGame(ctx context.Context, args []string) error {
	if len(args) < 0 {
		return errors.New("you must specify a game alias")
	}

	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	alias := args[0]
	if err := m.RemoveGame(alias); err != nil {
		fmt.Printf("error removing game: %v", err)
		return nil
	}

	if err := m.SaveConfig(); err != nil {
		fmt.Printf("error saving config: %v", err)
		return nil
	}

	fmt.Printf("game removed: %s", alias)
	return nil
}

/* List Games Command */
func listGames(ctx context.Context, _ []string) error {
	m, err := cmdutils.LoadManager(ctx)
	if err != nil {
		return err
	}

	if m.Config().Games.Len() < 1 {
		fmt.Printf("no games found.\n")
		return nil
	}

	fmt.Printf("%12s   %s\n", "ALIAS", "FOLDER")
	for alias, gameFolder := range m.Config().Games {
		fmt.Printf("%12s   %s\n", alias, gameFolder)
	}

	return nil
}
