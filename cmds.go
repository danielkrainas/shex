package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
)

var commandParser = flags.NewParser(nil, flags.None)

type executionContext struct {
	channels ChannelList
	homePath string
	config   *ManagerConfig
	profile  *Profile
	profiles map[string]*Profile
	args     []string
	remote   string
}

type appError struct {
	reason string
}

func (e appError) Error() string {
	return e.reason
}

type usageError struct{}

func (e usageError) Error() string {
	return ""
}

type logger func(format string, a ...interface{})

type commandExecutor func(current *executionContext, log logger) error

func consoleLogger(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func runInContext(fn commandExecutor) error {
	/*		var homePath string
			homeFlag := cmd.Flag.Lookup("home")
			if homeFlag != nil {
				homePath = homeFlag.Value.String()
			}*/
	homePath := ""
	config, err := loadManagerConfig(homePath)
	if err != nil {
		return err
	}

	profiles, err := getAvailableProfiles(config.ProfilesPath)
	if err != nil {
		return err
	}

	channels, err := loadAvailableChannels(config.ChannelsPath)
	if err != nil {
		return err
	}

	config.channels = channels
	if config.IncludeDefaultChannel {
		channels[defaultChannel.Alias] = defaultChannel
	}

	current := &executionContext{}
	current.homePath = homePath
	current.profiles = profiles
	current.config = config
	current.remote = getRemoteOrDefault(config.Remotes, config.ActiveRemote)
	current.profile = profiles[config.ActiveProfile]
	current.channels = channels
	return fn(current, consoleLogger)
}

func init() {
	commandParser.SubcommandsOptional = false
}

/*func loadCommandMap() map[string]*command {
	commands := map[string]*command{
		"stat": &command{
			action:    commandBootstrapper(execStat, 1),
			usageLine: "stat <path>",
			short:     "displays mod information",
			long:      "display mod information for a given mod package",
		},

		"pull": &command{
			action:    commandBootstrapper(execPull, 1),
			usageLine: "pull <name> [local_name]",
			short:     "pull a profile from a remote profile registry",
			long:      "Pulls a profile from a profile registry and imports it locally.",
		},

		"push": &command{
			action:    commandBootstrapper(execPush, 2),
			usageLine: "push <profile> <remote_profile>",
			short:     "pushes a profile to a profile registry",
			long:      "Pushes a profile to the profile registry under the name specified by remote_profile. If the remote profile already exists, it will publish a new version.",
		},
	}
}*/
