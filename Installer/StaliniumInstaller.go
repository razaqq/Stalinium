package main

import (
	"Stalinium/Installer/gui"
	"Stalinium/Installer/utils"
)

func main() {
	installs, _ := utils.GetWarshipsInstalls()
	ui := gui.StaliniumApp{Installs: installs}
	ui.Init()
}
