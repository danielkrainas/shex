package manager

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/danielkrainas/shex/api/client"
	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/fsutils"
)

func SyncProfile(p *v1.Profile) (int32, int32, error) {
	if p.Source == nil {
		return 0, 0, nil
	}

	rp, err := client.DownloadProfile(p.Source)
	if err != nil {
		return 0, 0, err
	}

	old := p.Revision
	p.Mods = rp.Mods
	p.Revision = rp.Revision
	return old, p.Revision, nil
}

func SaveProfileTo(p *v1.Profile, profilePath string) error {
	jsonContent, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(profilePath, jsonContent, 0777)
}

func SaveProfile(p *v1.Profile) error {
	if len(p.FilePath) < 1 {
		return errors.New("profile file path not set.")
	}

	return SaveProfileTo(p, p.FilePath)
}

func DropProfile(p *v1.Profile) error {
	if p.FilePath != "" && fsutils.FileExists(p.FilePath) {
		err := os.Remove(p.FilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadProfile(profilePath string) (v1.Profile, error) {
	var profile v1.Profile
	profile.FilePath = profilePath
	jsonContent, err := ioutil.ReadFile(profilePath)
	if err == nil {
		err = json.Unmarshal(jsonContent, &profile)
		if profile.Source.Type == v1.SOURCE_NONE {
			profile.Source = nil
		}
	}

	return profile, err
}

func pullProfile(source *v1.ProfileSource, localName string, profilesPath string) (*v1.Profile, error) {
	if source.Type != "remote" {
		return nil, errors.New("source type not supported")
	}

	profile, err := client.DownloadProfileAsLocal(source, localName)
	if err != nil {
		return nil, err
	}

	profile.FilePath = path.Join(profilesPath, localName+".json")
	return profile, SaveProfile(profile)
}

func pushProfile(profile *v1.Profile, remoteName string, endpoint string) (string, error) {
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

	res, err := client.PostContent(url, jsonContent)
	if err != nil {
		return "", err
	}

	return string(res[:]), nil
}

func createProfileSource(name string, location string) v1.ProfileSource {
	source := v1.ProfileSource{}
	source.Location = location
	source.Uid = name
	source.Type = "remote"
	return source
}

func LoadAllProfiles(profilesPath string) (map[string]*v1.Profile, error) {
	files, err := ioutil.ReadDir(profilesPath)
	result := make(map[string]*v1.Profile)
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
