package main

import ()

/* Sync Command */
type syncCommand struct {}

func (cmd *syncCommand) Execute(args []string) error {
    return usageError{}
}


/* Sync Profiles Command */
type syncProfilesCommand struct {}

func (cmd *syncProfilesCommand) Usage() error {
    
}

func (cmd *syncProfilesCommand) Execute(args []string) error {
    
}

func init() {
    sync, := commandParser.AddCommand("sync", "", "", &syncCommand{})
    sync.AddCommand("profiles", "", "", &syncProfilesCommand{})
}
