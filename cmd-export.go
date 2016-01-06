package main

import (
	"fmt"
)

/* Export Command */
type exportCommand struct{}

/* Export Profile Command */
type exportProfileCommand struct{}

func (cmd *exportProfileCommand) Usage() string {
	return "<path>"
}

func (cmd *exportProfileCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext, log logger) error {
		profileId := args[0]
		profile, ok := current.profiles[profileId]
		if !ok {
			return appError{fmt.Sprintf("[%s] not found\n", profileId)}
		}

		profilePath := args[1]
		err := saveProfile(profile, profilePath)
		if err != nil {
			return err
		}

		log("[%s] exported to: %s\n", profile.Id, profilePath)
		return nil
	})
}

func init() {
	export, _ := commandParser.AddCommand("export", "", "", &exportCommand{})

	export.AddCommand("profile", "exports a profile to a file", "Exports a profile to a file specified by path", &exportProfileCommand{})
}
