package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

const (
	ApiProfilesPath       = "profiles"
	ApiModsPath           = "mods"
	ApiModsMetaPathSuffix = "meta"
)

func downloadContents(url string) ([]byte, error) {
	resp, err := http.Get(url)
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return contents, err
	// smarter to io.Copy to local (temp?) file in case mods happen to be large
}

func downloadAndCacheContent(url string, cachePath string) (int64, error) {
	resp, err := http.Get(url)
	var read int64
	if err == nil {
		cacheFile, err := os.Create(cachePath)
		if err == nil {
			defer resp.Body.Close()
			defer cacheFile.Close()
			read, err = io.Copy(cacheFile, resp.Body)
		}
	}

	return read, err
}

func postContent(url string, bodyContent []byte) ([]byte, error) {
	body := bytes.NewBuffer(bodyContent)
	res, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func downloadModInfo(source string, mod *NameVersionToken) (*RemoteModInfo, error) {
	// NOTE: expect remote name to be something like user/package

	// http://somesite.com/mods/admin/cool-package/v/0.1.2/meta
	url := source + path.Join(ApiModsPath, mod.name, "v", mod.version, ApiModsMetaPathSuffix)
	contents, err := downloadContents(url)
	info := &RemoteModInfo{}
	info.source = url
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(contents, info)
	return info, err
}

func downloadModVersionList(source string, modName string) ([]string, error) {
	// http://somesite.com/mods/admin/cool-package/v
	url := source + path.Join(ApiModsPath, modName, "v")
	contents, err := downloadContents(url)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0)
	err = json.Unmarshal(contents, &versions)
	return versions, err
}

func downloadMod(source string, destPath string, info *RemoteModInfo) error {
	// http://somesite.com/mods/admin/cool-package/v/0.1.2
	url := source + path.Join(ApiModsPath, info.Name, "v", info.Version)

	_, err := downloadAndCacheContent(url, destPath)
	if err != nil {
		return err
	}

	return nil
}

func downloadProfileAsLocal(source *ProfileSource, localName string) (*Profile, error) {
	rp, err := downloadProfile(source)
	if err != nil {
		return nil, err
	}

	return makeLocalProfile(localName, rp), nil
}

func downloadProfile(source *ProfileSource) (*RemoteProfile, error) {
	url := path.Join(source.Location, ApiProfilesPath, source.Uid)
	jsonContent, err := downloadContents(url)
	if err != nil {
		return nil, err
	}

	remoteProfile := createRemoteProfile(source)
	err = json.Unmarshal(jsonContent, &remoteProfile)
	return remoteProfile, err
}
