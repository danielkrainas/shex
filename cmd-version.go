package main

import "fmt"

type versionCommand struct{}

func (cmd *versionCommand) Execute(args []string) error {
	return runInContext(func(current *executionContext) error {
		version := buildVersion
		if version == "" {
			version = "0.0.0-development"
		}

		fmt.Printf("%s %s\n", ProjectId, version)
		return nil
	})
}

func init() {
	commandParser.AddCommand("version", "", "", &versionCommand{})
}
