# goValoTUI

### This project is still under development.

goValoTUI is a terminal user interface used to display Valorant player's match statistics.

## Functionalities
- Basic CRUD interactions on a csv file to store player's info.
- Display up-to-date statistics about the player.

## Dependencies
- [tview](https://github.com/rivo/tview)
- [tcell](https://github.com/gdamore/tcell)
- [testify](https://github.com/stretchr/testify)
 
 ## How to
 #### Build & run :
 ```
 git clone https://github.com/Henrysoto/goValoTUI
 cd goValoTUI/
 go build
 ./goValoTUI
 ```
 #### Keybindings :
 - `Up/Down` arrow to navigate
 - `Ctrl+S` to remove player from list
 - `Enter` on input field will search for player tag
 - `Enter` on player list will fetch data for selected player

 ## Golang version
 `go1.20.2`
 
 ## Preview (WIP)
 ![goValoTUI](https://github.com/Henrysoto/goValoTUI/blob/master/screenshots/screenshot.png?raw=true)
