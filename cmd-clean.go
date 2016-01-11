package main

import (
	"path/filepath"
)

/* Clean Command */
type cleanCommand struct{}

func (cmd *cleanCommand) Usage() string {
	return "<cache>"
}

func (cmd *cleanCommand) Execute(args []string) error {
	if len(args) < 1 {
		return usageError{}
	}

	targetType := args[0]
	if targetType != "cache" {
		return usageError{}
	}

	return runInContext(func(current *executionContext, log logger) error {
		targetPath := filepath.Join(current.homePath, current.config.CachePath)
		if err := clearDirectory(targetPath); err != nil {
			return appError{"couldn't clear directory: " + targetPath}
		}

		log("cleared %s => %s", targetType, targetPath)
		return nil
	})
}

func init() {
	commandParser.AddCommand("clean", "remove groups of assets", "Removes all files from an asset directory", &cleanCommand{})
}
