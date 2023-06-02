package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	version = "0.0.1"
)

type TUI struct {
	app *tview.Application

	flexRoot         *tview.Flex
	flexLeft         *tview.Flex
	flexTopLeft      *tview.Flex
	flexTopLeftFlash *tview.Flex
	flexRight        *tview.Flex
	flexBotRight     *tview.Flex
	flexContainer    *tview.Flex

	displayMessage   *tview.TextView
	inputPlayerField *tview.InputField
	tablePlayers     *tview.Table
	tableMatch       *tview.Table
	formDetails      *tview.Form

	store *Store
}

func NewTUI(s *Store) (*TUI, error) {
	ui := TUI{
		app:   tview.NewApplication(),
		store: s,
	}

	ui.setupTUI()
	return &ui, nil
}

func (ui *TUI) setupTUI() {
	// flex grids
	ui.flexRoot = tview.NewFlex().SetDirection(tview.FlexRow)
	ui.flexContainer = tview.NewFlex()
	ui.flexLeft = tview.NewFlex().SetDirection(tview.FlexRow)
	ui.flexTopLeft = tview.NewFlex()
	ui.flexTopLeftFlash = tview.NewFlex()
	ui.flexRight = tview.NewFlex().SetDirection(tview.FlexRow)
	ui.flexBotRight = tview.NewFlex().SetDirection(tview.FlexRow)

	// components
	ui.displayMessage = tview.NewTextView()
	ui.inputPlayerField = tview.NewInputField()
	ui.tablePlayers = tview.NewTable()
	ui.tableMatch = tview.NewTable()
	ui.formDetails = tview.NewForm()

	// components config
	ui.flexRoot.
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(fmt.Sprintf(" goValoTUI %s ", version))
	ui.displayMessage.
		SetTextStyle(tcell.StyleDefault.Bold(true)).
		SetTextColor(tcell.ColorRed)
	ui.inputPlayerField.
		SetLabelStyle(tcell.StyleDefault.Bold(true)).
		SetBorder(true)
	ui.tableMatch.
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" Matches ")
	ui.tablePlayers.
		SetSelectable(true, false).
		SetFixed(1, 1).
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" Players ")
	ui.formDetails.
		SetItemPadding(0).
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(" Details ")

	// functionalities
	ui.setupEventHandlers()

	// layout
	ui.flexRoot.
		AddItem(ui.flexContainer, 0, 1, true)
	ui.flexContainer.
		AddItem(ui.flexLeft, 0, 4, true).
		AddItem(ui.flexRight, 0, 6, false)
	ui.flexLeft.
		AddItem(ui.flexTopLeftFlash, 1, 1, false).
		AddItem(ui.flexTopLeft, 3, 1, true).
		AddItem(ui.tablePlayers, 0, 1, false)
	ui.flexTopLeftFlash.
		AddItem(ui.displayMessage, 0, 1, false)
	ui.flexTopLeft.
		AddItem(ui.inputPlayerField, 0, 1, true)
	ui.flexRight.
		AddItem(ui.formDetails, 0, 1, false).
		AddItem(ui.flexBotRight, 0, 1, false)
	ui.flexBotRight.
		AddItem(ui.tableMatch, 0, 1, false)
}

func (ui *TUI) setupEventHandlers() {
	ui.inputPlayerField.SetDoneFunc(func(key tcell.Key) {
		ui.displayMessage.Clear()
		if key == tcell.KeyEnter {
			err := ui.inputSearch(ui.inputPlayerField.GetText())
			if err != nil {
				ui.displayError(err)
				return
			}
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyTAB {
			ui.app.SetFocus(ui.tablePlayers)
			return nil
		}
		return event
	})

	ui.tablePlayers.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ui.displayMessage.Clear()
		row, _ := ui.tablePlayers.GetSelection()
		switch event.Key() {
		case tcell.KeyEnter:
			player, err := ui.store.selectPlayerWithIndex(row) //selectPlayerWithName(ui.tablePlayers.GetCell(row, col).Text)
			if err != nil {
				ui.displayError(err)
				return event
			}
			ui.fillPlayerDetails(player)
			return event
		case tcell.KeyCtrlS:
			err := ui.store.removePlayerWithIndex(row)
			if err != nil {
				ui.displayError(err)
				return event
			}
			ui.tablePlayers.RemoveRow(row)
			return event
		case tcell.KeyUp:
			if row == 0 {
				ui.app.SetFocus(ui.inputPlayerField)
			}
			return event
		case tcell.KeyDown:
			if row+1 == ui.tablePlayers.GetRowCount() {
				ui.tablePlayers.Select(-1, 0)
			}
		}
		return event
	})
}

func (ui *TUI) fillPlayerTable(players []Player) {
	ui.tablePlayers.Clear()
	for index, player := range players {
		ui.tablePlayers.SetCellSimple(index, 0, GetUsername(player.Name, player.Tag))
	}
}

func (ui *TUI) fillPlayerDetails(player Player) {
	ui.formDetails.Clear(false)
	ui.displayMessage.Clear()
	ui.formDetails.SetTitle(fmt.Sprintf(" Player %s details ", player.Name))

	// Fetch account ranking data
	acc, err := GetData[AccountRank](fmt.Sprintf(BaseAPI+EndpointRank, player.Name, player.Tag))
	if err != nil {
		ui.displayError(err)
		return
	}

	var img image.Image
	if len(acc.Data) >= 1 {
		// Update player rank from store file if needed
		if player.Rank != acc.Data[0].Rank {
			player.Rank = acc.Data[0].Rank
			err = ui.store.updatePlayer(player)
			if err != nil {
				ui.displayError(err)
				return
			}
		}
		// Retrieve player card
		card, err := GetData[Account](fmt.Sprintf(BaseAPI+EndpointAccount, player.Name, player.Tag))
		if err != nil {
			ui.displayError(err)
			// Use default art encoded in base64
			base := base64.NewDecoder(base64.StdEncoding, strings.NewReader(artDefault64))
			img, err = png.Decode(base)
			if err != nil {
				ui.displayError(err)
				return
			}
		} else {
			// Otherwise fetch image from URL
			img, err = getImageFromURL(card.Data.Card.Wide)
			if err != nil {
				ui.displayError(fmt.Errorf("could not retrieve image from URL"))
			}
		}
	}

	// Fill form with data
	ui.formDetails.
		AddImage("Art", img, 0, 6, 256*256*256).
		AddTextView("", "", 1, 1, true, false).
		AddTextView("Username", GetUsername(player.Name, player.Tag), 0, 0, true, false).
		AddTextView("Rank", player.Rank, 0, 0, true, false).
		AddTextView("ELO", strconv.Itoa(acc.Data[0].Elo), 0, 0, true, false)

	// Player matches
	matches, err := GetData[MatchesData](fmt.Sprintf(BaseAPI+EndpointMatchesByName, player.Name, player.Tag))
	if err != nil {
		ui.displayError(err)
		return
	}

	// Update title with matches count
	ui.tableMatch.SetTitle(fmt.Sprintf(" Matches (%d) ", len(matches.Data)))

	// Set headers
	ui.tableMatch.SetCell(0, 0, tview.NewTableCell(FormatStringTable(" ", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 1, tview.NewTableCell(FormatStringTable("Map", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter)) //.SetSeparator(tview.BoxDrawingsLightDoubleDashVertical)
	ui.tableMatch.SetCell(0, 2, tview.NewTableCell(FormatStringTable("Agent", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 3, tview.NewTableCell(FormatStringTable("Score", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 4, tview.NewTableCell(FormatStringTable("Kills", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 5, tview.NewTableCell(FormatStringTable("Deaths", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 6, tview.NewTableCell(FormatStringTable("K/D Ratio", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 7, tview.NewTableCell(FormatStringTable("Headshots", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 8, tview.NewTableCell(FormatStringTable("Bodyshots", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.tableMatch.SetCell(0, 9, tview.NewTableCell(FormatStringTable("Damage", 12)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

	// Add empty row as separator
	ui.tableMatch.SetCellSimple(1, 0, " ")
	ui.tableMatch.SetCellSimple(1, 1, " ")
	ui.tableMatch.SetCellSimple(1, 2, " ")
	ui.tableMatch.SetCellSimple(1, 3, " ")
	ui.tableMatch.SetCellSimple(1, 4, " ")
	ui.tableMatch.SetCellSimple(1, 5, " ")
	ui.tableMatch.SetCellSimple(1, 6, " ")
	ui.tableMatch.SetCellSimple(1, 7, " ")
	ui.tableMatch.SetCellSimple(1, 8, " ")
	ui.tableMatch.SetCellSimple(1, 9, " ")

	// Set content
	var pID int
	for index, match := range matches.Data {
		// Get map played
		ui.tableMatch.SetCell(index+2, 1, tview.NewTableCell(match.Metadata.Map).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Get player index from match
		for i, p := range match.Players.AllPlayers {
			if p.Name == player.Name {
				pID = i
			}
		}

		agent := match.Players.AllPlayers[pID].Character
		score := match.Players.AllPlayers[pID].Stats.Score
		kills := match.Players.AllPlayers[pID].Stats.Kills
		deaths := match.Players.AllPlayers[pID].Stats.Deaths
		hs := match.Players.AllPlayers[pID].Stats.Headshots
		bs := match.Players.AllPlayers[pID].Stats.Bodyshots
		dmg := match.Players.AllPlayers[pID].DamageMade
		kd := float64(kills) / float64(deaths)

		// Get agent played
		ui.tableMatch.SetCell(index+2, 2, tview.NewTableCell(agent).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Get player score
		ui.tableMatch.SetCell(index+2, 3, tview.NewTableCell(fmt.Sprintf("%d", score)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Get player kills
		ui.tableMatch.SetCell(index+2, 4, tview.NewTableCell(fmt.Sprintf("%d", kills)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Get player deaths
		ui.tableMatch.SetCell(index+2, 5, tview.NewTableCell(fmt.Sprintf("%d", deaths)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Compute player kill/death ratio
		ui.tableMatch.SetCell(index+2, 6, tview.NewTableCell(fmt.Sprintf("%.3f", kd)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Get headshots count
		ui.tableMatch.SetCell(index+2, 7, tview.NewTableCell(fmt.Sprintf("%d", hs)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Get bodyshots count
		ui.tableMatch.SetCell(index+2, 8, tview.NewTableCell(fmt.Sprintf("%d", bs)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
		// Get damage count
		ui.tableMatch.SetCell(index+2, 9, tview.NewTableCell(fmt.Sprintf("%d", dmg)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	}
}

func (ui *TUI) inputSearch(fullPlayerName string) error {
	ui.inputPlayerField.SetText("")
	if !strings.Contains(fullPlayerName, "#") {
		return fmt.Errorf("username must contain a tag #")
	}

	testlength := strings.Split(fullPlayerName, "#")
	if len(testlength[1]) <= 1 || len(testlength[1]) > 5 {
		return fmt.Errorf("empty or invalid tag")
	}

	if len(fullPlayerName) < 7 || len(fullPlayerName) > 20 {
		return fmt.Errorf("username length is invalid")
	}

	// If player name already exists show details
	p1, err := ui.store.selectPlayerWithName(fullPlayerName)
	if p1 != (Player{}) {
		ui.displayError(fmt.Errorf("%v", p1))
		ui.fillPlayerDetails(p1)
		return nil
	}

	// Otherwise, fetch data from api and insert new player
	if err != nil {
		player := strings.Split(fullPlayerName, "#")
		newPlayer, err := GetData[AccountRank](fmt.Sprintf(BaseAPI+EndpointRank, player[0], player[1]))
		if err != nil {
			return err
		}

		var rank string
		if len(newPlayer.Data) <= 1 {
			rank = "Unranked"
		} else {
			rank = newPlayer.Data[0].Rank
		}

		players, err := ui.store.getPlayers()
		if err != nil {
			return err
		}
		players = append(players, Player{
			Name: player[0],
			Tag:  player[1],
			Rank: rank,
		})
		err = ui.store.insertPlayers(players)
		if err != nil {
			return err
		}

		// Update TUI
		ui.fillPlayerTable(players)
		ui.fillPlayerDetails(players[len(players)-1])
	}

	return nil
}

func (ui *TUI) displayError(err error) {
	ui.displayMessage.Clear()
	ui.displayMessage.SetText(fmt.Sprintf("[DEBUG]: %s", err.Error()))
}

// func initFormView(app *tview.Application, player Account) *tview.Flex {
// 	img, err := getImageFromURL(player.Data.Card.Small)
// 	if err != nil {
// 		log.Printf("Could not retrieve profile picture. Error: %v", err)
// 	}
// 	fw := tview.NewForm().
// 		AddImage(
// 			fmt.Sprintf("Level %d", player.Data.AccountLevel),
// 			img,
// 			0,
// 			0,
// 			256*256*256,
// 		).
// 		AddTextView("Player name", player.Data.Name, 40, 1, true, false).
// 		AddTextView("Player tag", fmt.Sprintf("#%s", player.Data.Tag), 20, 1, true, false).
// 		AddButton("Quit", func() {
// 			app.Stop()
// 		})
// 	fw.SetBorder(true).SetTitle("Valorant TUI Profile").SetTitleAlign(tview.AlignLeft)
// 	return tview.NewFlex().SetDirection(tview.FlexRow).AddItem(
// 		fw, 10, 1, true,
// 	)
// }
