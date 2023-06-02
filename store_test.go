package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var s *Store = NewTestStore()

func generateRandomPlayer() Player {
	return Player{
		Name: RandomPlayerName(),
		Tag:  RandomString(3),
		Rank: RandomPlayerRank(),
	}
}

func NewTestStore() *Store {
	return NewStore()
}

func TestGetPlayers(t *testing.T) {
	players, err := s.getPlayers()
	require.NoError(t, err)

	if l, _ := os.ReadFile(s.Filepath); len(l) > 1 {
		require.NotEmpty(t, players)
	} else {
		require.Empty(t, players)
	}
}

func TestInsertPlayers(t *testing.T) {
	players := []Player{
		generateRandomPlayer(),
		generateRandomPlayer(),
	}
	err := s.insertPlayers(players)
	require.NoError(t, err)
}

func TestUpdatePlayer(t *testing.T) {
	players, _ := s.getPlayers()
	player := players[0]
	player.Rank = RandomPlayerRank()
	err := s.updatePlayer(player)
	require.NoError(t, err)
}

func TestRemovePlayerWithIndex(t *testing.T) {
	players, err := s.getPlayers()
	require.NoError(t, err)
	require.NotEmpty(t, players)

	count := len(players)
	err = s.removePlayerWithIndex(int(RandomInt(0, int64(count))))
	require.NoError(t, err)

	players, err = s.getPlayers()
	require.NoError(t, err)
	require.NotEmpty(t, players)
	require.Equal(t, count-1, len(players))

	err = s.removePlayerWithIndex(len(players) + 2)
	require.Error(t, err)
}

func TestSelectPlayerWithIndex(t *testing.T) {
	players := []Player{generateRandomPlayer()}
	err := s.insertPlayers(players)
	require.NoError(t, err)

	players, err = s.getPlayers()
	require.NoError(t, err)

	player, err := s.selectPlayerWithIndex(int(RandomInt(0, int64(len(players)))))
	require.NoError(t, err)
	require.NotEmpty(t, player)

	player, err = s.selectPlayerWithIndex(len(players) + 2)
	require.Error(t, err)
	require.Empty(t, player)
}

func TestSelectPlayerWithName(t *testing.T) {
	players := []Player{generateRandomPlayer()}
	err := s.insertPlayers(players)
	require.NoError(t, err)

	players, err = s.getPlayers()
	require.NoError(t, err)

	p1 := players[int(RandomInt(0, int64(len(players))))]
	player, err := s.selectPlayerWithName(GetUsername(p1.Name, p1.Tag))
	require.NoError(t, err)
	require.NotEmpty(t, player)

	player, err = s.selectPlayerWithName(GetUsername(RandomPlayerName(), RandomString(4)))
	require.Error(t, err)
	require.Empty(t, player)
}
