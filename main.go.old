package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

/* App Error */
type appError struct {
	innerError error
	reason     string
}

func (e appError) Error() string {
	return e.reason
}

/* Usage Error */
type usageError struct{}

func (e usageError) Error() string {
	return ""
}

/* Execution Context */
type executionContext struct {
	channels ChannelMap
	homePath string
	config   *ManagerConfig
	profile  *Profile
	profiles map[string]*Profile
	channel  *Channel
}

/* App Options */
type AppOptions struct {
	HomePath string `long:"home" description:"override the home folder location"`
}

type commandExecutor func(current *executionContext) error

var (
	appOptions    = &AppOptions{}
	commandParser = flags.NewParser(appOptions, flags.None)
)

func runInContext(fn commandExecutor) error {
	config, err := loadManagerConfig(appOptions.HomePath)
	if err != nil {
		return err
	}

	profiles, err := loadAvailableProfiles(config.ProfilesPath)
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
	current.homePath = appOptions.HomePath
	current.profiles = profiles
	current.config = config
	if ch, ok := channels[config.ActiveRemote]; ok {
		current.channel = ch
	}

	current.profile = profiles[config.ActiveProfile]
	current.channels = channels
	return fn(current)
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

func main() {
	help := &helpCommand{}

	log.SetOutput(os.Stdout)
	_, err := commandParser.Parse()
	if err != nil {
		if aerr, ok := err.(appError); ok {
			fmt.Println(aerr.reason)
		} else if _, ok := err.(usageError); ok {
			help.Execute([]string{})
		} else if ferr, ok := err.(*flags.Error); ok {
			if ferr.Type == flags.ErrCommandRequired || ferr.Type == flags.ErrRequired {
				var args []string
				if commandParser.Active != nil {
					args = []string{commandParser.Active.Name}
				}

				help.Execute(args)
			}
		} else {
			println(err.Error())
			os.Exit(1)
		}
	}

	println("")
}
