package manager

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
)

type GameList map[string]string

func (list GameList) Detach(alias string) {
	delete(list, alias)
}

func (list GameList) Attach(alias string, gamePath string) {
	list[alias] = filepath.Clean(gamePath)
}

type Config struct {
	filePath              string
	ActiveProfile         string   `json:"active"`
	ActiveRemote          string   `json:"remote"`
	ProfilesPath          string   `json:"profiles"`
	ChannelsPath          string   `json:"channels"`
	IncludeDefaultChannel bool     `json:"includeDefaultChannel"`
	CachePath             string   `json:"cache"`
	Games                 GameList `json:"games"`
}

func SaveConfig(config *Config, homePath string) error {
	var err error
	if len(homePath) <= 0 {
		homePath, err = getDefaultHomePath()
		if err != nil {
			return err
		}
	}

	configPath := path.Join(homePath, HomeConfigName)
	jsonContent, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, jsonContent, 0777)
	if err != nil {
		return err
	}

	return nil
}

func NewConfig() *Config {
	config := &Config{}
	config.ActiveProfile = DefaultProfileName
	config.ActiveRemote = "default"
	config.Games = make(GameList)
	config.IncludeDefaultChannel = true
	return config
}

func LoadConfig(homePath string) (*Config, error) {
	config := &Config{}
	var err error
	if len(homePath) <= 0 {
		homePath, err = getDefaultHomePath()
		if err != nil {
			return nil, err
		}
	}

	configPath := filepath.Join(homePath, HomeConfigName)
	if err = ensureHomeDirectoryExists(homePath); err != nil {
		return nil, err
	}

	jsonContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonContent, config)
	if err == nil {
		config.filePath = configPath
		if len(config.ProfilesPath) < 1 {
			config.ProfilesPath = filepath.Join(homePath, HomeProfilesFolder)
		}

		if len(config.ChannelsPath) < 1 {
			config.ChannelsPath = filepath.Join(homePath, HomeChannelsFolder)
		}
	}

	return config, err
}
