package main

import (
	"fmt"
	"log"
)

/* Export Command */
type exportCommand struct{}

/* Export Profile Command */
type exportProfileCommand struct{}

func (cmd *exportProfileCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext) error {
		profileId := args[0]
		profile, ok := current.profiles[profileId]
		if !ok {
			return appError{nil, fmt.Sprintf("[%s] not found\n", profileId)}
		}

		profilePath := args[1]
		if err := profile.saveTo(profilePath); err != nil {
			return err
		}

		log.Printf("[%s] exported to: %s\n", profile.Id, profilePath)
		return nil
	})
}

func init() {
	export, _ := commandParser.AddCommand("export", "", "", &exportCommand{})
	export.AddCommand("profile", "", "", &exportProfileCommand{})
}
