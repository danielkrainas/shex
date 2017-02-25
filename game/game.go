package game

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

const DefaultGameManifestName = "shex.json"

const (
	ModsFolder = "mods"
)

type ModManifest struct {
	Info ModInfo `json:"info"`
}

type RemoteModInfo struct {
	source  string
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ModInfo struct {
	source     string
	Name       string `json:"name"`
	Version    int32  `json:"version"`
	SemVersion string `json:"semversion"`
}

func getModBasePath(modPath string) string {
	base := path.Base(modPath)
	return strings.TrimRight(base, ".smod")
}

func getZipResourceContent(zipPath string, resourcePath string, isRelative bool) ([]byte, error) {
	zip, err := zip.OpenReader(zipPath)
	if err == nil {
		defer zip.Close()
		for _, f := range zip.File {
			if isRelative {
				resourcePath = path.Join(strings.Split(f.Name, "/")[0], resourcePath)
				isRelative = false
			}

			if f.Name == resourcePath {
				r, err := f.Open()
				contents, err := ioutil.ReadAll(r)
				r.Close()
				return contents, err
			}
		}
	}

	return []byte{}, nil
}

func findVersion() (string, error) {
	versionRegEx := regexp.MustCompile("((?:Develop|(?:[a-zA-Z]+))\\-[0-9]+)")
	notesPath := "stonehearth/release_notes/release_notes.html"
	rawContent, err := getZipResourceContent("/home/daniel/Documents/stonehearth.smod", notesPath, false)
	result := ""
	if err == nil {
		if len(rawContent) <= 0 {
			return "", errors.New("could not find release notes")
		}

		notesContent := string(rawContent[:])
		for _, line := range strings.Split(notesContent, "\n") {
			if strings.Contains(line, "<h2>") && versionRegEx.MatchString(line) {
				result = versionRegEx.FindString(line)
				break
			}
		}
	}

	return result, err
}

func getModInfo(modPath string) (*ModInfo, error) {
	manifestPath := "/manifest.json"
	jsonContent, err := getZipResourceContent(modPath, manifestPath, true)
	if err != nil {
		return nil, err
	}

	if len(jsonContent) <= 0 {
		return nil, errors.New("could not find manifest file")
	}

	manifest := ModManifest{}
	err = json.Unmarshal(jsonContent, &manifest)
	if err == nil {
		if len(manifest.Info.SemVersion) <= 0 {
			manifest.Info.SemVersion = fmt.Sprintf("%d.0.0", manifest.Info.Version)
		}

		manifest.Info.source = modPath
	}

	return &manifest.Info, err
}
