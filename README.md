# goValoTUI

### This project is still under development.

goValoTUI is a terminal user interface used to display Valorant player's match statistics.

## Functionalities
- Basic CRUD operations on a csv file to store player's info.
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
 - `Up/Down` arrow to navigate.
 - `Ctrl+S` to remove player from players list.
 - `Enter` on input field will search for the player's profile and add it to the list.
 - `Enter` on player list will fetch data for selected player.

 ## Golang version
 `go1.20.2`
 
 ## Preview (WIP)
 ![goValoTUI](https://github.com/Henrysoto/goValoTUI/blob/master/screenshots/screenshot.png?raw=true)

## TODO
- Edit form player details to display properly art picture and item's padding.
- Retrieve all matches and not only last five matches.
  - Add pagination to match's table.
- Add dropdown menu to choose between Competitive and Unrated matches.
- Add columns sorting.
