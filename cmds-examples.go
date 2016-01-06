// +build !windows

package main

var commandExamples = map[string]string{
	"goble":     "$ goble install foo/mod",
	"stat":      "$ goble stat <path>",
	"version":   "$ goble version",
	"install":   "$ goble install user/mod-name@version",
	"mods":      "$ goble mods \nor a specific game: \n$ goble mods game2",
	"uninstall": "$ goble uninstall user/mod-name",
	"pull":      "$ goble pull remote/name localname",
	"push":      "$ goble push default foouser/default",
	"export":    "$ goble export default /path/to/export/file.json",
	"create":    "$ goble create demo \nor with a specified path: \n$ goble create demo /path/to/save/file.json",
	"drop":      "$ goble drop default",
	"import":    "$ goble import /path/to/config.json",
	"use":       "$ goble use default",
	"profiles":  "$ goble profiles",
	"games":     "$ goble games",
	"detach":    "$ goble detach default",
	"attach":    "$ goble attach default /path/to/stonehearth",
}
