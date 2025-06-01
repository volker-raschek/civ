package main

import "git.cryptic.systems/volker.raschek/civ/cmd"

var version string

func main() {
	_ = cmd.Execute(version)
}
