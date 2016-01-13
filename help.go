package main

type helpData struct {
	short string
	long  string
	usage string
}

var helpIndex = map[string]*helpData{
	"import": &helpData{
		long:  "Imports an asset into the manager",
		short: "import and asset",
	},

	"import profile": &helpData{
		long:  "Imports a file as a profile config.",
		short: "imports a profile config file",
		usage: "<path>",
	},

	"export": &helpData{
		long:  "Removes all files from an asset directory.",
		short: "remove groups of assets",
		usage: "<path>",
	},

	"export profile": &helpData{
		long:  "Exports a profile to a file specified by path.",
		short: "exports a profile to a file",
		usage: "<path>",
	},

	"clean": &helpData{
		long:  "Removes all files from an asset directory.",
		short: "remove groups of assets",
		usage: "<cache>",
	},

	"clean cache": &helpData{
		long:  "Clears the cache folder of all contents.",
		short: "empties the cache folder",
	},

	"add": &helpData{
		long:  "Adds assets to the manager using one of the available <command>.",
		usage: "",
	},

	"add profile": &helpData{
		long:  "Creates a new mod profile with the specified id. If a path argument is supplied, the profile won't be imported and will be saved to the path specified.",
		short: "add a profile config",
		usage: "<id> [path]",
	},

	"add game": &helpData{
		short: "add a game folder",
		long:  "Adds the game folder at the specified location to the manager. <alias> may be omitted and \"default\" will be assumed.",
		usage: "<alias> <game_path>",
	},

	"add channel": &helpData{
		long:  "Adds a remote channel to the manager.",
		short: "add a remote channel",
		usage: "<alias> <endpoint> [options]",
	},

	"help": &helpData{
		long:  "Display help and usage information.",
		short: "displays help information",
		usage: "<term> [<sub term 1>...<sub term N>]",
	},

	"install": &helpData{
		long:  "Install a mod into the current profile",
		short: "install a mod",
		usage: "self|[<mod1>...<modN>]",
	},

	"install self": &helpData{
		long:  "Installs goble on the current system.",
		short: "install goble",
		usage: "[path]",
	},

	"list": &helpData{
		long:  "List manager data.",
		short: "list manager data",
	},

	"list mods": &helpData{
		long:  "Lists the mods installed in the default or specified game.",
		short: "lists the mods that are installed",
	},

	"list games": &helpData{
		long:  "Lists the games currently attached to the manager.",
		short: "lists the game folders attached",
	},

	"list profiles": &helpData{
		long:  "List the available mod profiles.",
		short: "lists available profiles",
		usage: "",
	},

	"list config": &helpData{
		long:  "Lists the current config settings.",
		short: "lists config settings",
		usage: "",
	},

	"list channels": &helpData{
		long:  "Lists channels available in the manager.",
		short: "list available channels",
		usage: "",
	},

	"version": &helpData{
		long:  "Prints the manager version.",
		short: "prints the manager version",
		usage: "",
	},

	"set": &helpData{
		long:  "Changes a manager config setting.",
		short: "change a config setting",
		usage: "<key> <value>",
	},

	"uninstall": &helpData{
		long:  "Uninstall a mod from the current profile.",
		short: "uninstall a mod",
		usage: "<mod> [<mod1>...<modN>]",
	},

	"uninstall self": &helpData{
		long:  "Uninstalls goble from the system.",
		short: "uninstall goble",
		usage: "self|[<mod1>...<modN>]",
	},

	"use": &helpData{
		long:  "sets the active mod profile",
		short: "sets the active profile",
		usage: "<profile>",
	},

	"remove": &helpData{
		long:  "Removes a manager asset.",
		short: "remove manager assets",
	},

	"remove profile": &helpData{
		long:  "Removes the profile from the manager.",
		short: "remove a profile",
		usage: "<profile>",
	},

	"remove game": &helpData{
		long:  "Removes a game by the specified alias from the manager.",
		short: "remove a game folder",
		usage: "<game>",
	},

	"remove channel": &helpData{
		long:  "Removes a channel from the manager.",
		short: "remove a channel",
		usage: "<channel>",
	},
}

func init() {
	/* Add aliases references */
	helpIndex["i"] = helpIndex["install"]
	helpIndex["i self"] = helpIndex["install self"]

	helpIndex["u"] = helpIndex["uninstall"]
	helpIndex["u self"] = helpIndex["uninstall self"]
}
