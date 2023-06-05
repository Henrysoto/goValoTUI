package main

import (
	"fmt"
	"log"
)

var players []Player

func main() {
	fmt.Printf("[goValoTUI %s]\n", version)

	s := NewStore()
	players, err := s.getPlayers()
	if err != nil {
		log.Panicln(err)
	}

	ui := NewTUI(s)
	ui.fillPlayerTable(players)

	if err = ui.app.SetRoot(ui.flexRoot, true).EnableMouse(true).Run(); err != nil {
		log.Panicln(err)
	}
}
