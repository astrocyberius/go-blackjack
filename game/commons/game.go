package commons

import (
	"github.com/astrocyberius/go-blackjack/config"
	ds "github.com/astrocyberius/go-blackjack/datastructs"
	"github.com/astrocyberius/go-blackjack/ui"
)

func InitGame(deck *ds.Deck52, gameRound *ds.GameRoundData, gameState ds.GameState, playerName string) {
	*gameRound = ds.GameRoundData{PlayingCards: GetPlayingCards(ShuffleCardDeck(deck)),
		Player: ds.Player{Name: playerName, Money: config.PlayerMoney}}

	InitGameRound(gameRound, gameState)
	ui.PrintWelcomePlayer(gameRound.Player.Name)
}

func InitGameRound(gameRound *ds.GameRoundData, gameState ds.GameState) {
	player := &gameRound.Player

	player.RightHand = []ds.PlayingCard{}
	player.RightHandPlayingDouble = false
	player.BetRightHand = 0

	player.LeftHand = []ds.PlayingCard{}
	player.LeftHandPlayingDouble = false
	player.BetLeftHand = 0

	player.MakeRightHandActive()
	player.Split = false

	gameRound.DealerCards = []ds.PlayingCard{}

	gameRound.GameState = gameState
}
