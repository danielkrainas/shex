package manager

var (
	defaultChannel = &Channel{
		Alias:    "default",
		Endpoint: "127.0.0.1:6231/",
		Protocol: "http",
	}
)

type ChannelMap map[string]*Channel

type Channel struct {
	filePath string
	Alias    string `json:"alias"`
	Protocol string `json:"protocol"`
	Endpoint string `json:"endpoint"`
}

func (ch *Channel) saveTo(channelPath string) error {
	jsonContent, err := json.Marshal(ch)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(channelPath, jsonContent, 0777)
}

func (ch *Channel) save() error {
	return ch.saveTo(ch.filePath)
}

func (ch *Channel) remove() error {
	return os.Remove(ch.filePath)
}

func loadChannel(channelPath string) (*Channel, error) {
	channel := &Channel{
		filePath: channelPath,
	}

	jsonContent, err := ioutil.ReadFile(channelPath)
	if err == nil {
		err = json.Unmarshal(jsonContent, channel)
	}

	return channel, err
}

func LoadAllChannels(pathDir string) ([]*Channel, error) {

}

func loadAvailableChannels(channelsPath string) (ChannelMap, error) {
	files, err := ioutil.ReadDir(channelsPath)
	result := make(ChannelMap)
	if err == nil {
		for _, f := range files {
			isJson, err := filepath.Match("*.json", f.Name())
			if err != nil {
				return nil, err
			}

			if isJson {
				channel, err := loadChannel(filepath.Join(channelsPath, f.Name()))
				if err != nil {
					return nil, err
				}

				result[channel.Alias] = channel
			}
		}
	}

	return result, err
}