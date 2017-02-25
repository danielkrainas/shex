package game

import (
	"strings"

	"github.com/danielkrainas/shex/configuration"
)

type ModList map[string]string

type GameManifest struct {
	Mods ModList `json:"mods"`
}

type NameVersionToken struct {
	name    string
	version string
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

func uninstallMod(config *configuration.ManagerConfig, gamePath string, profile *Profile, name string) (*ModInfo, error) {
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

func installMod(config *configuration.ManagerConfig, gamePath string, profile *Profile, modToken *NameVersionToken) (*ModInfo, error) {
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
