package main

type helpData struct {
	short string
	long  string
	usage string
}

var helpIndex = map[string]*helpData{
	"import": {
		long:  "Imports an asset into the manager",
		short: "import and asset",
	},

	"import profile": {
		long:  "Imports a file as a profile config.",
		short: "imports a profile config file",
		usage: "<path>",
	},

	"export": {
		long:  "Removes all files from an asset directory.",
		short: "remove groups of assets",
		usage: "<path>",
	},

	"export profile": {
		long:  "Exports a profile to a file specified by path.",
		short: "exports a profile to a file",
		usage: "<path>",
	},

	"clean": {
		long:  "Removes all files from an asset directory.",
		short: "remove groups of assets",
		usage: "<cache>",
	},

	"clean cache": {
		long:  "Clears the cache folder of all contents.",
		short: "empties the cache folder",
	},

	"add": {
		long:  "Adds assets to the manager using one of the available <command>.",
		usage: "",
	},

	"add profile": {
		long:  "Creates a new mod profile with the specified id. If a path argument is supplied, the profile won't be imported and will be saved to the path specified.",
		short: "add a profile config",
		usage: "<id> [path]",
	},

	"add game": {
		short: "add a game folder",
		long:  "Adds the game folder at the specified location to the manager. <alias> may be omitted and \"default\" will be assumed.",
		usage: "<alias> <game_path>",
	},

	"add channel": {
		long:  "Adds a remote channel to the manager.",
		short: "add a remote channel",
		usage: "<alias> <endpoint> [options]",
	},

	"help": {
		long:  "Display help and usage information.",
		short: "displays help information",
		usage: "<term> [<sub term 1>...<sub term N>]",
	},

	"install": {
		long:  "Install a mod into the current profile",
		short: "install a mod",
		usage: "self|[<mod1>...<modN>]",
	},

	"install self": {
		long:  "Installs shex on the current system.",
		short: "install shex",
		usage: "[path]",
	},

	"list": {
		long:  "List manager data.",
		short: "list manager data",
	},

	"list mods": {
		long:  "Lists the mods installed in the default or specified game.",
		short: "lists the mods that are installed",
	},

	"list games": {
		long:  "Lists the games currently attached to the manager.",
		short: "lists the game folders attached",
	},

	"list profiles": {
		long:  "List the available mod profiles.",
		short: "lists available profiles",
		usage: "",
	},

	"list config": {
		long:  "Lists the current config settings.",
		short: "lists config settings",
		usage: "",
	},

	"list channels": {
		long:  "Lists channels available in the manager.",
		short: "list available channels",
		usage: "",
	},

	"version": {
		long:  "Prints the manager version.",
		short: "prints the manager version",
		usage: "",
	},

	"set": {
		long:  "Changes a manager config setting.",
		short: "change a config setting",
		usage: "<key> <value>",
	},

	"uninstall": {
		long:  "Uninstall a mod from the current profile.",
		short: "uninstall a mod",
		usage: "<mod> [<mod1>...<modN>]",
	},

	"uninstall self": {
		long:  "Uninstalls shex from the system.",
		short: "uninstall shex",
		usage: "self|[<mod1>...<modN>]",
	},

	"use": {
		long:  "sets the active mod profile",
		short: "sets the active profile",
		usage: "<profile>",
	},

	"remove": {
		long:  "Removes a manager asset.",
		short: "remove manager assets",
	},

	"remove profile": {
		long:  "Removes the profile from the manager.",
		short: "remove a profile",
		usage: "<profile>",
	},

	"remove game": {
		long:  "Removes a game by the specified alias from the manager.",
		short: "remove a game folder",
		usage: "<game>",
	},

	"remove channel": {
		long:  "Removes a channel from the manager.",
		short: "remove a channel",
		usage: "<channel>",
	},

	"sync": {
		long:  "Sync a manager asset.",
		short: "sync manager assets with their sources",
	},

	"sync profiles": {
		long:  "Syncs all applicable profiles with their valid remote sources and updates them to the latest revision.",
		short: "sync and update all remote profiles",
	},

	"sync profile": {
		long:  "Syncs a local profile with its remote source.",
		short: "sync a remote profile",
		usage: "<profile>",
	},
}

func init() {
	/* Add aliases references */
	helpIndex["i"] = helpIndex["install"]
	helpIndex["i self"] = helpIndex["install self"]

	helpIndex["u"] = helpIndex["uninstall"]
	helpIndex["u self"] = helpIndex["uninstall self"]
}

/*func loadCommandMap() map[string]*command {
	commands := map[string]*command{
		"stat": &command{
			action:    commandBootstrapper(execStat, 1),
			usageLine: "stat <path>",
			short:     "displays mod information",
			long:      "display mod information for a given mod package",
		},

		"pull": &command{
			action:    commandBootstrapper(execPull, 1),
			usageLine: "pull <name> [local_name]",
			short:     "pull a profile from a remote profile registry",
			long:      "Pulls a profile from a profile registry and imports it locally.",
		},

		"push": &command{
			action:    commandBootstrapper(execPush, 2),
			usageLine: "push <profile> <remote_profile>",
			short:     "pushes a profile to a profile registry",
			long:      "Pushes a profile to the profile registry under the name specified by remote_profile. If the remote profile already exists, it will publish a new version.",
		},
	}
}*/
