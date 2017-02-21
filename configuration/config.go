package configuration

type GameList map[string]string

type ManagerConfig struct {
	filePath              string
	ActiveProfile         string   `json:"active"`
	ActiveRemote          string   `json:"remote"`
	ProfilesPath          string   `json:"profiles"`
	ChannelsPath          string   `json:"channels"`
	IncludeDefaultChannel bool     `json:"includeDefaultChannel"`
	CachePath             string   `json:"cache"`
	Games                 GameList `json:"games"`
}
