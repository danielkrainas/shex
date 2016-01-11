package main

import (
	"log"
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

	return runInContext(func(current *executionContext) error {
		targetPath := filepath.Join(current.homePath, current.config.CachePath)
		if err := clearDirectory(targetPath); err != nil {
			return appError{err, "couldn't clear directory: " + targetPath}
		}

		log.Printf("cleared %s => %s", targetType, targetPath)
		return nil
	})
}

func init() {
	commandParser.AddCommand("clean", "remove groups of assets", "Removes all files from an asset directory", &cleanCommand{})
}
