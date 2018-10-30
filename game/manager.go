package game

import (
	"github.com/astrocyberius/go-blackjack/config"
	ds "github.com/astrocyberius/go-blackjack/datastructs"
	"github.com/astrocyberius/go-blackjack/game/commons"
	"github.com/astrocyberius/go-blackjack/game/event"
	"github.com/astrocyberius/go-blackjack/game/state"
	"github.com/astrocyberius/go-blackjack/ui"
	"os"
)

var (
	deck      *ds.Deck52
	gameRound ds.GameRoundData
)

// InitGame initializes the game state
func InitGame() {
	event.RegisterEventHandler(event.GameMenuEvent, handleGameMenuEvent)
	event.RegisterEventHandler(event.GameInfoEvent, handleGameInfoEvent)
	event.RegisterEventHandler(event.GameRoundEvent, handleGameRoundEvent)
	event.RegisterEventHandler(event.GameQuitEvent, handleGameQuitEvent)
	event.RegisterEventHandler(event.PlayerGameHitEvent, handlePlayerGameHitEvent)
	event.RegisterEventHandler(event.PlayerGameStandEvent, handlePlayerGameStandEvent)

	ui.PrintIntro()
	deck = commons.CreateCardDeck()
	gameState := &state.PlayerWithSingleHandState{GameRoundData: &gameRound, Deck: deck}
	commons.InitGame(deck, &gameRound, gameState, ui.GetPlayerName())
}

// StartGameLoop starts the game loop
func StartGameLoop() {
	event.AddEvent(&event.Event{EventType: event.GameMenuEvent})
	event.HandleEvents()
}

func handleGameMenuEvent() {
	ui.PrintMenuAndHandlePlayerInput(&ui.MenuOptionCallback{Play: playGameRoundCallback, Info: infoCallback,
		Quit: quitCallback})
}

func playGameRoundCallback() {
	event.AddEvent(&event.Event{EventType: event.GameRoundEvent})
}

func infoCallback() {
	event.AddEvent(&event.Event{EventType: event.GameInfoEvent})
}

func quitCallback() {
	event.AddEvent(&event.Event{EventType: event.GameQuitEvent})
}

func handleGameInfoEvent() {
	ui.PrintInfo()
	event.AddEvent(&event.Event{EventType: event.GameMenuEvent})
}

func handleGameQuitEvent() {
	ui.PrintQuit()
	os.Exit(0)
}

func handleGameRoundEvent() {
	player := &gameRound.Player

	quitToMainMenu := false
	player.BetRightHand, quitToMainMenu = ui.GetPlayerBetInput(GetBetOptions, player.Money, config.MinimumBetValue, config.MaximumBetValue)
	if quitToMainMenu {
		gameState := &state.PlayerWithSingleHandState{GameRoundData: &gameRound, Deck: deck}
		commons.InitGame(deck, &gameRound, gameState, player.Name)
		event.AddEvent(&event.Event{EventType: event.GameMenuEvent})
		return
	}

	player.Money -= player.BetRightHand
	player.RightHand = append(player.RightHand, commons.PopCardFromDeck(deck, &gameRound.PlayingCards))
	ui.PrintPlayerHitCard((player.RightHand)[len(player.RightHand)-1])

	gameRound.DealerCards = append(gameRound.DealerCards, commons.PopCardFromDeck(deck, &gameRound.PlayingCards))
	ui.PrintDealerHitCard(gameRound.DealerCards[len(gameRound.DealerCards)-1])

	gameRound.GameState.Hit()
}

func handlePlayerGameHitEvent() {
	gameRound.GameState.Hit()
}

func handlePlayerGameStandEvent() {
	gameRound.GameState.Stand()
}
