package mods

import (
	"strings"

	"github.com/danielkrainas/shex/api/v1"
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
