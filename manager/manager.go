package manager

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
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
