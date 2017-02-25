package manager

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path"
	"path/filepath"
)

const (
	DefaultGameManifestName = "shex.json"
	DefaultGameName         = "default"
	DefaultProfileName      = "default"
	HomeConfigName          = "config.json"
	HomeProfilesFolder      = "profiles"
	HomeCacheFolder         = "cache"
	HomeChannelsFolder      = "channels"
	defaultHomeFolder       = ".shex"
)

func getGameOrDefault(games GameList, name string) string {
	if name == "" {
		name = DefaultGameName
	}

	gamePath, ok := games[name]
	if !ok {
		return ""
	}

	return gamePath
}

func getDefaultHomePath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(u.HomeDir, string(filepath.Separator)+defaultHomeFolder), nil
}

func saveManagerConfig(config *ManagerConfig, homePath string) error {
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

func createManagerConfig() *ManagerConfig {
	config := &ManagerConfig{}
	config.ActiveProfile = DefaultProfileName
	config.ActiveRemote = "default"
	config.Games = make(GameList)
	config.IncludeDefaultChannel = true
	return config
}

func loadManagerConfig(homePath string) (*ManagerConfig, error) {
	config := &ManagerConfig{}
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
