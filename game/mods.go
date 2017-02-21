package game

type ModList map[string]string

type GameManifest struct {
	Mods ModList `json:"mods"`
}

type NameVersionToken struct {
	name    string
	version string
}
