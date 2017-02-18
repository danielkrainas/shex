package main

var commandExamples = map[string]string{
	"shex":      "shex install foo/mod",
	"stat":      "shex stat <path>",
	"version":   "shex version",
	"install":   "shex install user/mod-name@version",
	"mods":      "shex mods \nor a specific game: \nshex mods game2",
	"uninstall": "shex uninstall user/mod-name",
	"pull":      "shex pull remote/name localname",
	"push":      "shex push default foouser/default",
	"export":    "shex export default " + `"c:\path\to\export\file.json"`,
	"create":    "shex create demo \nor with a specified path: \nshex create demo " + `"c:\path\to\save\file.json"`,
	"drop":      "shex drop default",
	"import":    "shex import " + `"c:\path\to\config.json"`,
	"use":       "shex use default",
	"profiles":  "shex profiles",
	"games":     "shex games",
	"detach":    "shex detach game1",
	"attach":    "shex attach default " + `"c:\program files\steam\steamapps\common\stonehearth"`,
}
