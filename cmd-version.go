package main

import "log"

type versionCommand struct{}

func (cmd *versionCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext) error {
		version := buildVersion
		if version == "" {
			version = "0.0.0-development"
		}

		log.Printf("%s %s\n", ProjectId, version)
		return nil
	})
}

func init() {
	commandParser.AddCommand("version", "prints the manager version", "Prints the manager version.", &versionCommand{})
}
