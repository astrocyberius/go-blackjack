package config

import ds "github.com/astrocyberius/go-blackjack/datastructs"

const (
	// DealerStandsAtOrAboveHandValue dealer stands when the hand value is at or above the specified value
	DealerStandsAtOrAboveHandValue = 17

	// PlayerMoney the total amount of money in the player's wallet
	PlayerMoney ds.TMoney = 2500

	// MinimumBetValue the minimum bet value
	MinimumBetValue ds.TMoney = 5

	// MaximumBetValue the maximum bet value
	MaximumBetValue ds.TMoney = 10000
)
