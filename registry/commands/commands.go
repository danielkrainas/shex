package commands

import (
	"github.com/danielkrainas/shex/api/v1"
)

type DeleteMod struct {
	Name string
}

type StoreMod struct {
	New bool
	Mod *v1.ModInfo
}
