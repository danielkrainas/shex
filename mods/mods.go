package mods

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/fsutils"
)

type ModManifest struct {
	Info v1.ModInfo `json:"info"`
}

type GameManifest struct {
	Mods v1.ModList `json:"mods"`
}

func GetLocalModPathName(remoteName string, version string) string {
	parts := strings.Split(remoteName, "/")
	return parts[0] + "_" + parts[1] + "-" + version + ".smod"
}

/*func isVersionOkay(version string, versionConstraint string) bool {

}
*/

func CreateGameManifest() *GameManifest {
	manifest := &GameManifest{
		Mods: make(v1.ModList),
	}

	return manifest
}

func LoadGameManifest(gamePath string) (*GameManifest, error) {
	manifest := CreateGameManifest()
	manifestPath := path.Join(gamePath, DefaultGameManifestName)
	if !fsutils.FileExists(manifestPath) {
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

func SaveGameManifest(gamePath string, manifest *GameManifest) error {
	jsonContent, err := json.Marshal(manifest)
	if err != nil {
		return err
	}

	manifestPath := path.Join(gamePath, DefaultGameManifestName)
	return ioutil.WriteFile(manifestPath, jsonContent, 0777)
}
