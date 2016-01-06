package main

var commandExamples = map[string]string{
	"goble":     "goble install foo/mod",
	"stat":      "goble stat <path>",
	"version":   "goble version",
	"install":   "goble install user/mod-name@version",
	"mods":      "goble mods \nor a specific game: \ngoble mods game2",
	"uninstall": "goble uninstall user/mod-name",
	"pull":      "goble pull remote/name localname",
	"push":      "goble push default foouser/default",
	"export":    "goble export default " + `"c:\path\to\export\file.json"`,
	"create":    "goble create demo \nor with a specified path: \ngoble create demo " + `"c:\path\to\save\file.json"`,
	"drop":      "goble drop default",
	"import":    "goble import " + `"c:\path\to\config.json"`,
	"use":       "goble use default",
	"profiles":  "goble profiles",
	"games":     "goble games",
	"detach":    "goble detach game1",
	"attach":    "goble attach default " + `"c:\program files\steam\steamapps\common\stonehearth"`,
}
