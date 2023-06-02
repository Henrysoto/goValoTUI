package main

import (
	"fmt"
	"log"
)

var players []Player

func main() {
	fmt.Println("[goValoTUI]")

	// account, err := getAccount()
	// if err != nil {
	// 	log.Panicln(err)
	// }

	s := NewStore()
	players, err := s.getPlayers()
	if err != nil {
		log.Panicln(err)
	}

	ui, err := NewTUI(s)
	if err != nil {
		log.Panicln(err)
	}

	ui.fillPlayerTable(players)
	ui.displayError(fmt.Errorf("test message"))

	if err = ui.app.SetRoot(ui.flexRoot, true).EnableMouse(true).Run(); err != nil {
		log.Panicln(err)
	}
}
