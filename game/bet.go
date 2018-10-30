package game

import (
	ds "github.com/astrocyberius/go-blackjack/datastructs"
)

var gameChips = [...]ds.TMoney{5, 10, 25, 50, 100, 500, 1000, 5000}

// GetBetOptions returns the possible bet options (chips to play with) calculated from the provided money
func GetBetOptions(playerMoney ds.TMoney) []ds.TMoney {
	var betOptions []ds.TMoney

	for _, v := range gameChips {
		if playerMoney < v {
			break
		}
		betOptions = append(betOptions, v)
	}

	return betOptions
}
