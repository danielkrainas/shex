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

func ensureHomeDirectoryExists(homePath string) error {
	if !dirExists(homePath) {
		err := os.Mkdir(homePath, 0777)
		if err != nil {
			return err
		}
	}

	configPath := filepath.Join(homePath, HomeConfigName)
	profilesPath := filepath.Join(homePath, HomeProfilesFolder)
	cachePath := filepath.Join(homePath, HomeCacheFolder)
	channelsPath := filepath.Join(homePath, HomeChannelsFolder)
	if !fileExists(configPath) {
		defaultConfig := createManagerConfig()
		defaultConfig.ProfilesPath = profilesPath

		jsonContent, err := json.Marshal(&defaultConfig)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(configPath, jsonContent, 0777); err != nil {
			return err
		}
	}

	if !dirExists(cachePath) {
		err := os.Mkdir(cachePath, 0777)
		if err != nil {
			return err
		}
	}

	if !dirExists(channelsPath) {
		if err := os.Mkdir(channelsPath, 0777); err != nil {
			return err
		}
	}

	if !dirExists(profilesPath) {
		err := os.Mkdir(profilesPath, 0777)
		if err != nil {
			return err
		}
	}

	defaultProfilePath := path.Join(profilesPath, DefaultProfileName+".json")
	if !fileExists(defaultProfilePath) {
		defaultProfile := Profile{}
		defaultProfile.Id = DefaultProfileName
		defaultProfile.Mods = make(map[string]string)
		defaultProfile.Name = strings.ToTitle(DefaultProfileName)

		jsonContent, err := json.Marshal(&defaultProfile)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(defaultProfilePath, jsonContent, 0777); err != nil {
			return err
		}
	}

	return nil
}
