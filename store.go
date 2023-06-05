package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Store struct {
	Filepath string
}

type Player struct {
	Name string
	Tag  string
	Rank string
}

// Create dirs and store file if it does not exist
func createFile(dir string, filepath string) error {
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}
	if _, err := os.ReadFile(filepath); err != nil {
		fp, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer fp.Close()
	}
	return nil
}

// Open or create a new store file
func NewStore() *Store {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err)
	}
	dirPath := filepath.Join(home, ".config", "goValoTUI", "data")
	filePath := filepath.Join(dirPath, "players.csv")
	err = createFile(dirPath, filePath)
	if err != nil {
		log.Panicln(err)
	}
	return &Store{Filepath: filePath}
}

// Parse store file and return Player struct array
func (s *Store) getPlayers() ([]Player, error) {
	var players []Player
	file, err := os.ReadFile(s.Filepath)
	if err != nil {
		return players, err
	}

	for _, line := range strings.Split(string(file), "\n") {
		items := strings.Split(line, ",")
		if len(items) < 3 {
			continue
		}
		players = append(players, Player{
			Name: items[0],
			Tag:  items[1],
			Rank: items[2],
		})
	}

	return players, nil
}

// Insert new player data to store file
func (s *Store) insertPlayers(players []Player) error {
	db, err := s.getPlayers()
	if err != nil {
		return err
	}

	content, err := s.getData()
	if err != nil {
		return err
	}

	// flag
	exists := false

	for _, player := range players {
		for _, dbItem := range db {
			if dbItem.Name == player.Name && dbItem.Tag == player.Tag {
				exists = true
			}
		}
		if !exists {
			content = append(content, fmt.Sprintf("%s,%s,%s",
				player.Name,
				player.Tag,
				player.Rank,
			))
		}
		exists = false
	}
	if len(content) > 1 {
		err = s.writeData(strings.Join(content, "\n"))
		if err != nil {
			log.Printf("Could not save data. Error: %v\n", err)
			return err
		}
	}

	return nil
}

// Update store line containing player
func (s *Store) updatePlayer(player Player) error {
	lines, err := s.getData()
	if err != nil {
		return err
	}

	for index, line := range lines {
		if strings.Contains(line, fmt.Sprintf("%s,%s", player.Name, player.Tag)) {
			lines[index] = fmt.Sprintf("%s,%s,%s", player.Name, player.Tag, player.Rank)
		}
	}

	err = s.writeData(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}

	return nil
}

// Remove player at index from store file
func (s *Store) removePlayerWithIndex(index int) error {
	lines, err := s.getData()
	if err != nil {
		return err
	}

	if index <= len(lines) && index >= 0 {
		lines = append(lines[:index], lines[index+1:]...)
		err = s.writeData(strings.Join(lines, "\n"))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no data found at given index '%d'", index)
	}

	return nil
}

// Parse store file and return string array containing each lines
func (s *Store) getData() ([]string, error) {
	var data []string
	fp, err := os.ReadFile(s.Filepath)
	if err != nil {
		return data, err
	}

	data = strings.Split(string(fp), "\n")
	for i, item := range data {
		data[i] = strings.Trim(item, "\n")
		if len(item) < 1 || item == "" {
			data = append(data[:i], data[i+1:]...)
		}
	}
	return data, nil
}

// Write string data to store file (csv format)
func (s *Store) writeData(data string) error {
	if len(data) < 1 {
		return fmt.Errorf("no data to write")
	}

	err := os.WriteFile(s.Filepath, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Return player from store file at given index
func (s *Store) selectPlayerWithIndex(index int) (Player, error) {
	var player Player
	data, err := s.getData()
	if err != nil {
		return player, err
	}

	if len(data) > 1 {
		for i, line := range data {
			if i == index {
				item := strings.Split(line, ",")
				player.Name = item[0]
				player.Tag = item[1]
				player.Rank = item[2]
			}
		}
	} else {
		return player, fmt.Errorf("no data to select")
	}

	if player != (Player{}) {
		return player, nil
	} else {
		return player, fmt.Errorf("could not find player at index %d", index)
	}
}

// Return player from store file with given full username string "username#tag"
func (s *Store) selectPlayerWithName(username string) (Player, error) {
	var player Player
	data, err := s.getData()
	if err != nil {
		return player, err
	}

	if len(data) > 1 {
		for _, line := range data {
			item := strings.Split(line, ",")
			fullusername := GetUsername(item[0], item[1])
			if fullusername == username {
				player.Name = item[0]
				player.Tag = item[1]
				player.Rank = item[2]
			}
		}
	} else {
		return player, fmt.Errorf("no data to select")
	}

	if player != (Player{}) {
		return player, nil
	} else {
		return player, fmt.Errorf("could not find player '%s'", username)
	}
}
