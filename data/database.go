package data

import (
	"fmt"
	"math/rand"
	"time"
)

// Game is a game
type Game struct {
	ID             string        `json:"id"`
	HidingSpots    []*HidingSpot `json:"hidingSpots"`
	TurnsRemaining int           `json:"turnsRemainins"`
}

// HidingSpot is
type HidingSpot struct {
	ID             string `json:"id"`
	HasBeenChecked bool   `json:"hasBeenChecked"`
	HasTreasure    bool   `json:"hasTreasure"`
}

var game *Game

// Init initializes our singular game
func Init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	winningSpotID := int(r.Int31n(9))
	hidingSpots := make([]*HidingSpot, 9)
	for i := 0; i < 9; i++ {
		winningSpot := i == winningSpotID
		hidingSpots[i] = &HidingSpot{
			ID:             fmt.Sprintf("%d", i),
			HasBeenChecked: false,
			HasTreasure:    winningSpot,
		}
	}

	game = &Game{string(1), hidingSpots, 3}
}

// CheckHidingSpotForTreasure will mark a hiding spot checked and decrement the turns remaining
func (g *Game) CheckHidingSpotForTreasure(id string) {
	for _, v := range g.HidingSpots {
		if v.HasTreasure && v.HasBeenChecked {
			return
		}
	}
	hidingSpot := g.GetHidingSpot(id)
	if hidingSpot.HasBeenChecked {
		return
	}
	g.TurnsRemaining--
	hidingSpot.HasBeenChecked = true
}

// GetGame returns the gameboard
func GetGame() *Game {
	return game
}

// GetHidingSpot returns a hiding spot on a game
func (g *Game) GetHidingSpot(id string) *HidingSpot {
	for _, v := range g.HidingSpots {
		if v.ID == id {
			return v
		}
	}
	return nil
}

// GetHidingSpots returns the collection of hiding spots
func (g *Game) GetHidingSpots() []*HidingSpot {
	return g.HidingSpots
}

// GetTurnsRemaining returns the number of guesses left available
func (g *Game) GetTurnsRemaining() int {
	return g.TurnsRemaining
}
