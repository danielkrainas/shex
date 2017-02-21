package main

import (
	"encoding/json"
	"errors"
	//"github.com/hashicorp/go-version"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

type ActionError string

func copyFile(src string, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return -1, err
	}

	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return -1, err
	}

	read, err := io.Copy(dstFile, srcFile)
	if err != nil {
		dstFile.Close()
		return -1, err
	}

	err = dstFile.Close()
	return read, err
}

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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return true
}

func dirExists(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil || !stat.IsDir() {
		return false
	}

	return true
}

func clearDirectory(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		p := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			err = os.RemoveAll(p)
		} else {
			err = os.Remove(p)
		}

		if err != nil {
			return err
		}
	}

	return nil
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

func parseNameVersionToken(pair string) *NameVersionToken {
	token := &NameVersionToken{}
	parts := strings.Split(pair, "@")
	token.name = parts[0]
	if len(parts) > 1 {
		token.version = parts[1]
	} else {
		token.version = "latest"
	}

	return token
}

func getLocalModPathName(remoteName string, version string) string {
	parts := strings.Split(remoteName, "/")
	return parts[0] + "_" + parts[1] + "-" + version + ".smod"
}

/*func isVersionOkay(version string, versionConstraint string) bool {

}

func isModCached(config *ManagerConfig) bool {

}*/

func uninstallMod(config *ManagerConfig, gamePath string, profile *Profile, name string) (*ModInfo, error) {
	mod := &ModInfo{}
	gameManifest, err := loadGameManifest(gamePath)
	if err != nil {
		return mod, err
	}

	delete(profile.Mods, name)
	err = profile.save()
	if err != nil {
		return mod, err
	}

	modsPath := filepath.Join(gamePath, GameModsFolder)
	version, ok := gameManifest.Mods[name]
	if ok {
		modPath := filepath.Join(modsPath, getLocalModPathName(name, version))
		if fileExists(modPath) {
			err = os.Remove(modPath)
			if err != nil {
				return mod, err
			}
		}

		delete(gameManifest.Mods, name)
		err = saveGameManifest(gamePath, gameManifest)
		if err != nil {
			return mod, err
		}
	} else {
		//fmt.Printf("not installed in %s: \"%s\"\n", modsPath, name)
	}

	return mod, err
}

func installMod(config *ManagerConfig, gamePath string, profile *Profile, modToken *NameVersionToken) (*ModInfo, error) {
	mod := &ModInfo{}
	ch, ok := config.channels[config.ActiveRemote]
	if !ok {
		return mod, errors.New("channel not found: " + config.ActiveRemote)
	}

	source := ch.Protocol + "://" + ch.Endpoint
	remoteInfo, err := downloadModInfo(source, modToken)
	if err != nil {
		return mod, err
	}

	localName := getLocalModPathName(remoteInfo.Name, remoteInfo.Version)
	localPath := filepath.Join(gamePath, GameModsFolder, localName)
	err = downloadMod(source, localPath, remoteInfo)
	if err != nil {
		return mod, err
	}

	gameManifest, err := loadGameManifest(gamePath)
	if err != nil {
		return mod, err
	}

	profileVersion := modToken.version
	if profileVersion != "latest" {
		profileVersion = "^" + profileVersion
	}

	profile.Mods[remoteInfo.Name] = profileVersion
	err = profile.save()
	if err != nil {
		return mod, err
	}

	gameManifest.Mods[remoteInfo.Name] = remoteInfo.Version
	err = saveGameManifest(gamePath, gameManifest)
	if err != nil {
		return mod, err
	}

	return getModInfo(localPath)
}

func createGameManifest() *GameManifest {
	manifest := &GameManifest{}
	manifest.Mods = make(ModList)
	return manifest
}

func loadGameManifest(gamePath string) (*GameManifest, error) {
	manifest := createGameManifest()
	manifestPath := path.Join(gamePath, DefaultGameManifestName)
	if !fileExists(manifestPath) {
		return manifest, nil
	}

	jsonContent, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return manifest, err
	}

	err = json.Unmarshal(jsonContent, manifest)
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}

func saveGameManifest(gamePath string, manifest *GameManifest) error {
	jsonContent, err := json.Marshal(manifest)
	if err != nil {
		return err
	}

	manifestPath := path.Join(gamePath, DefaultGameManifestName)
	return ioutil.WriteFile(manifestPath, jsonContent, 0777)
}

/*func execStat(current *executionContext) error {
	modPath := args[0]
	info, err := getModInfo(modPath)
	if err != nil {
		return appError{err, "Could not find mod information"}
	}

	log.Printf("[%s]\n", modPath)
	log.Printf("name: %s\nversion: %d\nsem version: %s\n", info.Name, info.Version, info.SemVersion)
	return nil
}

func execPull(current *executionContext) error {
	remoteName := args[0]
	localName := path.Base(remoteName)
	if len(current.args) > 1 {
		localName = args[1]
	}

	if _, ok := current.profiles[localName]; ok {
		return appError{nil, fmt.Sprintf("[%s] already exists", localName)}
	}

	var ok bool
	// TODO: put this together as its own part later
	current.remote = getDefaultRemote()
	if current.config.ActiveRemote != DefaultRemoteName {
		current.remote, ok = current.config.Remotes[current.config.ActiveRemote]
		if !ok {
			return appError{nil, fmt.Sprintf("remote \"%s\" not found\n", current.config.ActiveRemote)}
		}
	}

	source := createProfileSource(remoteName, current.remote)
	profile, err := pullProfile(&source, localName, current.config.ProfilesPath)
	if err != nil {
		return appError{err, "Could not pull profile from the server"}
	}

	log.Printf("pulled [%s] to: %s\n", profile.Source.Uid, profile.filePath)
	return nil
}

func execPush(current *executionContext) error {
	profileId := args[0]
	profile, ok := current.profiles[profileId]
	if !ok {
		return appError{nil, fmt.Sprintf("[%s] not found\n", profileId)}
	}

	remoteName := args[1]
	if remoteName != current.config.ActiveRemote {
		current.remote = getRemoteOrDefault(current.config.Remotes, remoteName)
	}

	// TODO: pull out core logic into own func or something
	remote := getDefaultRemote()
	if current.config.ActiveRemote != "default" {
		remote, ok = current.config.Remotes[current.config.ActiveRemote]
		if !ok {
			return appError{nil, fmt.Sprintf("remote \"%s\" not found\n", current.config.ActiveRemote)}
		}
	}

	version, err := pushProfile(profile, remoteName, remote)
	if err != nil {
		return appError{err, "Could not push profile to server"}
	}

	log.Printf("[%s] pushed to %s as %s@%s\n", profileId, current.config.ActiveRemote, remoteName, version)
	log.Printf("import with: `shex pull %s@%s`\n", remoteName, version)
	return nil
}*/
