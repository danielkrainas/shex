package main

import ()

/* Set Command */
type setCommand struct{}

func (cmd *setCommand) Execute(args []string) error {
	if len(args) < 2 {
		return usageError{}
	}

	target := args[0]
	value := args[1]
	return runInContext(func(current *executionContext) error {
		switch target {
		case "profile":
			if _, ok := current.profiles[value]; ok {
				current.config.ActiveProfile = value
			} else {
				return appError{nil, "profile not found"}
			}

			break

		case "channel":
			if _, ok := current.channels[value]; ok {
				current.config.ActiveRemote = value
			} else {
				return appError{nil, "channel not found"}
			}

			break

			/*case "game":
			if game, ok := current.config.Games[value]; ok {
				current.config.ActiveGame = value
			} else {
				return appError{"game not found"}
			}

			break*/
		default:
			return usageError{}
		}

		return saveManagerConfig(current.config, current.homePath)
	})
}

func init() {
	commandParser.AddCommand("set", "", "", &setCommand{})
}
