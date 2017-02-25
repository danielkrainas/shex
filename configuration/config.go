package configuration

import (
	"path/filepath"
)

type GameList map[string]string

func (list GameList) Detach(alias string) {
	delete(list, alias)
}

func (list GameList) Attach(alias string, gamePath string) {
	list[alias] = filepath.Clean(gamePath)
}

type ManagerConfig struct {
	filePath              string
	ActiveProfile         string   `json:"active"`
	ActiveRemote          string   `json:"remote"`
	ProfilesPath          string   `json:"profiles"`
	ChannelsPath          string   `json:"channels"`
	IncludeDefaultChannel bool     `json:"includeDefaultChannel"`
	CachePath             string   `json:"cache"`
	Games                 GameList `json:"games"`
}
