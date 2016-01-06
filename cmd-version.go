package main

type versionCommand struct{}

func (cmd *versionCommand) Execute(args []string) error {
	return runInContext(execVersion)
}

func execVersion(current *executionContext, log logger) error {
	version := buildVersion
	if version == "" {
		version = "0.0.0-development"
	}

	log("%s %s\n", ProjectId, version)
	return nil
}

func init() {
	commandParser.AddCommand("version", "prints the manager version", "Prints the manager version.", &versionCommand{})
}
