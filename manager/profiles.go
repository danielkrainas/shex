package manager

import (
	"github.com/danielkrainas/shex/game"
)

const (
	SOURCE_REMOTE = "remote"
	SOURCE_NONE   = ""
)

type ProfileSource struct {
	Type     string `json:"type"`
	Uid      string `json:"uid"`
	Location string `json:"url"`
}

type RemoteProfile struct {
	source   *ProfileSource
	Name     string  `json:"name"`
	Mods     ModList `json:"mods"`
	Revision int32   `json:"rev"`
}

type Profile struct {
	filePath string
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	Mods     game.ModList   `json:"mods"`
	Source   *ProfileSource `json:"source"`
	Revision int32          `json:"rev"`
}

func (p *Profile) sync() (int32, int32, error) {
	if p.Source == nil {
		return 0, 0, nil
	}

	rp, err := downloadProfile(p.Source)
	if err != nil {
		return 0, 0, err
	}

	old := p.Revision
	p.Mods = rp.Mods
	p.Revision = rp.Revision
	return old, p.Revision, nil
}

func (p *Profile) contains(modName string) bool {
	_, ok := p.Mods[modName]
	return ok
}

func (p *Profile) saveTo(profilePath string) error {
	jsonContent, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(profilePath, jsonContent, 0777)
}

func (p *Profile) save() error {
	if len(p.filePath) < 1 {
		return errors.New("profile file path not set.")
	}

	return p.saveTo(p.filePath)
}

func (p *Profile) drop() error {
	if p.filePath != "" && fileExists(p.filePath) {
		err := os.Remove(p.filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func createProfile(id string) *Profile {
	profile := &Profile{}
	profile.Id = id
	profile.Name = strings.Title(id)
	profile.Mods = make(ModList)
	profile.Revision = 1
	profile.Source = nil
	return profile
}

func loadProfile(profilePath string) (Profile, error) {
	var profile Profile
	profile.filePath = profilePath
	jsonContent, err := ioutil.ReadFile(profilePath)
	if err == nil {
		err = json.Unmarshal(jsonContent, &profile)
		if profile.Source.Type == SOURCE_NONE {
			profile.Source = nil
		}
	}

	return profile, err
}

func createRemoteProfile(source *ProfileSource) *RemoteProfile {
	profile := &RemoteProfile{}
	profile.Mods = make(ModList)
	profile.source = source
	return profile
}

func makeLocalProfile(localName string, remote *RemoteProfile) *Profile {
	profile := createProfile(localName)
	profile.Source = remote.source
	profile.Revision = remote.Revision
	profile.Mods = remote.Mods
	return profile
}

func pullProfile(source *ProfileSource, localName string, profilesPath string) (*Profile, error) {
	if source.Type != "remote" {
		return nil, errors.New("source type not supported")
	}

	profile, err := downloadProfileAsLocal(source, localName)
	if err != nil {
		return nil, err
	}

	profile.filePath = path.Join(profilesPath, localName+".json")
	return profile, profile.save()
}

func pushProfile(profile *Profile, remoteName string, endpoint string) (string, error) {
	url := endpoint + "profiles/" + remoteName
	remoteProfile := *profile
	remoteProfile.Id = remoteName
	remoteProfile.Name = strings.Title(path.Base(remoteName))
	remoteProfile.Source = nil
	remoteProfile.Revision = 0
	jsonContent, err := json.Marshal(remoteProfile)
	if err != nil {
		return "", err
	}

	res, err := postContent(url, jsonContent)
	if err != nil {
		return "", err
	}

	return string(res[:]), nil
}

func createProfileSource(name string, location string) ProfileSource {
	source := ProfileSource{}
	source.Location = location
	source.Uid = name
	source.Type = "remote"
	return source
}

func loadAvailableProfiles(profilesPath string) (map[string]*Profile, error) {
	files, err := ioutil.ReadDir(profilesPath)
	result := make(map[string]*Profile)
	if err == nil {
		for _, f := range files {
			isJson, err := filepath.Match("*.json", f.Name())
			if err != nil {
				return nil, err
			}

			if isJson {
				profile, err := loadProfile(path.Join(profilesPath, f.Name()))
				if err != nil {
					return nil, err
				}

				result[profile.Id] = &profile
			}
		}
	}

	return result, err
}

func LoadAllProfiles(pathDir string) ([]*Profile, error) {

}
