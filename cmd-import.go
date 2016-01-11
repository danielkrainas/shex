package main

import (
	"log"
	"path/filepath"
	"strings"
)

/* Import Command */
type importCommand struct{}

/* Import Profile Command */
type importProfileCommand struct{}

func (cmd *importProfileCommand) Usage() string {
	return "<path>"
}

func (cmd *importProfileCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	return runInContext(func(current *executionContext) error {
		profilePath := args[0]
		profile, err := loadProfile(profilePath)
		if err != nil {
			return appError{err, "Could not load profile import"}
		}

		if len(args) > 1 {
			profile.Id = strings.ToLower(args[1])
			profile.Name = strings.Title(args[1])
		}

		newProfilePath := filepath.Join(current.config.ProfilesPath, profile.Id+".json")
		if err := saveProfile(&profile, newProfilePath); err != nil {
			return err
		}

		log.Printf("imported \"%s\" to: %s\n", profile.Id, newProfilePath)
		return nil
	})
}

func init() {
	imp, _ := commandParser.AddCommand("import", "import and asset", "import an asset into the manager", &importCommand{})

	imp.AddCommand("profile", "imports a profile config file", "imports a file as a profile config", &importProfileCommand{})
}
