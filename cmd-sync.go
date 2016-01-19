package main

import (
	"fmt"
)

func reportSyncResult(artifactName string, fromVersion string, toVersion string) {
	if fromVersion == toVersion {
		fmt.Printf("%s %-20s\n", artifactName, "OK")
	} else {
		fmt.Printf("%s %-20s->%s\n", artifactName, fromVersion, toVersion)
	}
}

func reportProfileSyncResult(p *Profile, from int32, to int32) {
	if from == to {
		fmt.Printf("%s @%d => no updates available\n", p.Name, from)
	} else {
		fmt.Printf("%s @%d => @%d\n", p.Name, from, to)
	}
}

/* Sync Command */
type syncCommand struct{}

func (cmd *syncCommand) Execute(args []string) error {
	return usageError{}
}

/* Sync Profiles Command */
type syncProfilesCommand struct{}

func (cmd *syncProfilesCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext) error {
		for _, p := range current.profiles {
			if p.Source == nil {
				continue
			}

			from, to, err := p.sync()
			if err != nil {
				return appError{err, "couldn't sync with remote server."}
			}

			err = saveProfile(p, p.filePath)
			if err != nil {
				return appError{err, "couldn't save profile."}
			}

			reportProfileSyncResult(p, from, to)
		}

		return nil
	})
}

/* Sync Profile Command */
type syncProfileCommand struct{}

func (cmd *syncProfileCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext) error {
		profile := current.profile
		if len(args) > 0 {
			var ok bool
			profile, ok = current.profiles[args[0]]
			if !ok {
				return appError{nil, "profile not found."}
			}
		}

		if profile.Source == nil {
			return appError{nil, "not a remote profile."}
		}

		from, to, err := profile.sync()
		if err != nil {
			return appError{err, "couldn't sync with remote server."}
		}

		err = saveProfile(profile, profile.filePath)
		if err != nil {
			return appError{err, "couldn't save profile."}
		}

		reportProfileSyncResult(profile, from, to)
		return nil
	})
}

func init() {
	sync, _ := commandParser.AddCommand("sync", "", "", &syncCommand{})
	sync.AddCommand("profiles", "", "", &syncProfilesCommand{})
	sync.AddCommand("profile", "", "", &syncProfileCommand{})
}
