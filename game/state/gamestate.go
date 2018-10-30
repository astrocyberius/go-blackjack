package state

import (
	"github.com/astrocyberius/go-blackjack/config"
	ds "github.com/astrocyberius/go-blackjack/datastructs"
	"github.com/astrocyberius/go-blackjack/game/commons"
	"github.com/astrocyberius/go-blackjack/game/event"
	"github.com/astrocyberius/go-blackjack/ui"
)

// PlayerWithSingleHandState player without split
type PlayerWithSingleHandState struct {
	GameRoundData *ds.GameRoundData
	Deck          *ds.Deck52
}

func (playerWithSingleHandState *PlayerWithSingleHandState) Hit() {
	gameRoundData := playerWithSingleHandState.GameRoundData

	player := &gameRoundData.Player
	playerCards := &player.RightHand
	handPlayingDouble := player.RightHandPlayingDouble

	// give a new card to the player
	*playerCards = append(*playerCards, commons.PopCardFromDeck(playerWithSingleHandState.Deck, &gameRoundData.PlayingCards))
	ui.PrintPlayerHitCard((*playerCards)[len(*playerCards)-1])

	playerHandValue := commons.GetHandValue(*playerCards)

	if playerHandValue < commons.BlackjackWinPoints && !handPlayingDouble {
		withDoubleOption := isHandEligibleForDouble(*playerCards, player.Money, player.BetRightHand)
		withSplitOption := isHandEligibleForSplit(*playerCards, player.Money, player.BetRightHand)

		callbacks := &ui.PlayerGameOptionCallback{Hit: playerGameHitCallback, Double: playerGameDoubleCallback,
			Stand: playerGameStandCallback, Split: playerGameSplitCallback}
		ui.PrintGameOptionsAndHandlePlayerInput(callbacks, playerWithSingleHandState.Deck, gameRoundData, player,
			ui.CardsWithValue{Cards: *playerCards, Value: playerHandValue},
			ui.CardsWithValue{Cards: gameRoundData.DealerCards, Value: commons.GetHandValue(gameRoundData.DealerCards)},
			withDoubleOption, withSplitOption)
	} else {
		gameRoundData.GameState.Stand()
	}
}

func (playerWithSingleHandState *PlayerWithSingleHandState) Double() {
	gameRoundData := playerWithSingleHandState.GameRoundData

	player := &gameRoundData.Player
	player.Money -= player.BetRightHand
	player.BetRightHand *= 2
	player.RightHandPlayingDouble = true

	event.AddEvent(&event.Event{EventType: event.PlayerGameHitEvent})
}

func (playerWithSingleHandState *PlayerWithSingleHandState) Stand() {
	gameRoundData := playerWithSingleHandState.GameRoundData

	player := &gameRoundData.Player
	handleRoundOutcome(playerWithSingleHandState.Deck, gameRoundData.PlayingCards, player.RightHand,
		&player.Money, &player.BetRightHand, &gameRoundData.DealerCards)

	gameState := &PlayerWithSingleHandState{GameRoundData: gameRoundData, Deck: playerWithSingleHandState.Deck}
	if player.Money < config.MinimumBetValue {
		ui.PrintNotEnoughMoneyLeftToPlay()
		commons.InitGame(playerWithSingleHandState.Deck, gameRoundData, gameState, player.Name)
		event.AddEvent(&event.Event{EventType: event.GameMenuEvent})
	} else {
		commons.InitGameRound(gameRoundData, gameState)
		event.AddEvent(&event.Event{EventType: event.GameRoundEvent})
	}
}

// PlayerWithDoubleHandState player with split
type PlayerWithDoubleHandState struct {
	GameRoundData *ds.GameRoundData
	Deck          *ds.Deck52
}

func (playerWithDoubleHandState *PlayerWithDoubleHandState) Hit() {
	gameRoundData := playerWithDoubleHandState.GameRoundData
	player := &gameRoundData.Player

	var playerCards *[]ds.PlayingCard
	var handPlayingDouble bool
	var bet ds.TMoney
	if player.IsRightHandActive() {
		playerCards = &player.RightHand
		handPlayingDouble = player.RightHandPlayingDouble
		bet = player.BetRightHand
	} else {
		playerCards = &player.LeftHand
		handPlayingDouble = player.LeftHandPlayingDouble
		bet = player.BetLeftHand
	}

	// give a new card to the player
	*playerCards = append(*playerCards, commons.PopCardFromDeck(playerWithDoubleHandState.Deck, &gameRoundData.PlayingCards))
	ui.PrintPlayerHitCard((*playerCards)[len(*playerCards)-1])

	playerHandValue := commons.GetHandValue(*playerCards)

	if playerHandValue < commons.BlackjackWinPoints && !handPlayingDouble {
		withDoubleOption := isHandEligibleForDouble(*playerCards, player.Money, bet)

		callbacks := &ui.PlayerGameOptionCallback{Hit: playerGameHitCallback, Double: playerGameDoubleCallback,
			Stand: playerGameStandCallback, Split: nil}
		ui.PrintGameOptionsAndHandlePlayerInput(callbacks, playerWithDoubleHandState.Deck, gameRoundData, player,
			ui.CardsWithValue{Cards: *playerCards, Value: playerHandValue},
			ui.CardsWithValue{Cards: gameRoundData.DealerCards, Value: commons.GetHandValue(gameRoundData.DealerCards)},
			withDoubleOption, false)
	} else if player.IsRightHandActive() {
		ui.PrintSwitchingToPlayerLeftHand()
		player.MakeLeftHandActive()
		event.AddEvent(&event.Event{EventType: event.PlayerGameHitEvent})
	} else {
		gameRoundData.GameState.Stand()
	}
}

func (playerWithDoubleHandState *PlayerWithDoubleHandState) Double() {
	gameRoundData := playerWithDoubleHandState.GameRoundData
	player := &gameRoundData.Player

	if player.IsRightHandActive() {
		player.RightHandPlayingDouble = true
		player.Money -= player.BetRightHand
		player.BetRightHand *= 2
	} else {
		player.LeftHandPlayingDouble = true
		player.Money -= player.BetLeftHand
		player.BetLeftHand *= 2
	}

	event.AddEvent(&event.Event{EventType: event.PlayerGameHitEvent})
}

func (playerWithDoubleHandState *PlayerWithDoubleHandState) Stand() {
	gameRoundData := playerWithDoubleHandState.GameRoundData
	player := &gameRoundData.Player

	if player.IsRightHandActive() {
		ui.PrintSwitchingToPlayerLeftHand()
		player.MakeLeftHandActive()
		event.AddEvent(&event.Event{EventType: event.PlayerGameHitEvent})
		return
	}

	ui.PrintMessage("Outcome for right hand cards", true, true)
	handleRoundOutcome(playerWithDoubleHandState.Deck, gameRoundData.PlayingCards, player.RightHand, &player.Money,
		&player.BetRightHand, &gameRoundData.DealerCards)

	ui.PrintMessage("Outcome for left hand cards", true, true)
	handleRoundOutcome(playerWithDoubleHandState.Deck, gameRoundData.PlayingCards, player.LeftHand, &player.Money,
		&player.BetLeftHand, &gameRoundData.DealerCards)

	gameState := &PlayerWithSingleHandState{GameRoundData: gameRoundData, Deck: playerWithDoubleHandState.Deck}
	if player.Money < config.MinimumBetValue {
		ui.PrintNotEnoughMoneyLeftToPlay()
		commons.InitGame(playerWithDoubleHandState.Deck, gameRoundData, gameState, player.Name)
		event.AddEvent(&event.Event{EventType: event.GameMenuEvent})
	} else {
		commons.InitGameRound(gameRoundData, gameState)
		event.AddEvent(&event.Event{EventType: event.GameRoundEvent})
	}
}

func isHandEligibleForDouble(hand []ds.PlayingCard, money ds.TMoney, bet ds.TMoney) bool {
	return len(hand) == 2 && money >= bet
}

func isHandEligibleForSplit(hand []ds.PlayingCard, money ds.TMoney, bet ds.TMoney) bool {
	return len(hand) == 2 && hand[0].Rank == hand[1].Rank && money >= bet
}

func handleRoundOutcome(deck *ds.Deck52, playingCards []ds.PlayingCard, playerCards []ds.PlayingCard,
	playerMoney *ds.TMoney, playerBet *ds.TMoney, dealerCards *[]ds.PlayingCard) {

	playerHandValue := commons.GetHandValue(playerCards)
	dealerHandValue := commons.GetHandValue(*dealerCards)

	if playerHandValue <= commons.BlackjackWinPoints && dealerHandValue < config.DealerStandsAtOrAboveHandValue {
		// draw the second card for the dealer
		*dealerCards = append(*dealerCards, commons.PopCardFromDeck(deck, &playingCards))
		ui.PrintDealerHitCard((*dealerCards)[len(*dealerCards)-1])

		// draw more cards for the dealer if necessary
		for commons.GetHandValue(*dealerCards) < config.DealerStandsAtOrAboveHandValue {
			*dealerCards = append(*dealerCards, commons.PopCardFromDeck(deck, &playingCards))
			ui.PrintDealerHitCard((*dealerCards)[len(*dealerCards)-1])
		}
	}

	dealerHandValue = commons.GetHandValue(*dealerCards)

	drawGame := false
	playerWins := false
	playerHasBlackjackHand := commons.IsBlackjackHand(playerCards, playerHandValue)
	dealerHasBlackjackHand := commons.IsBlackjackHand(*dealerCards, dealerHandValue)

	if playerHasBlackjackHand {
		if dealerHasBlackjackHand {
			drawGame = true
		} else {
			playerWins = true
		}
	} else if playerHandValue <= commons.BlackjackWinPoints {
		playerWins = dealerHandValue > commons.BlackjackWinPoints || playerHandValue > dealerHandValue
		drawGame = playerHandValue == dealerHandValue
	}

	moneyWonOrLost := *playerBet
	if drawGame {
		*playerMoney += moneyWonOrLost
	} else if playerWins {
		moneyWonOrLost = *playerBet * 2
		if playerHasBlackjackHand {
			moneyWonOrLost += *playerBet * 0.5
		}
		*playerMoney += moneyWonOrLost
	}

	ui.PrintRoundOutcome(playerCards, playerHandValue, *playerMoney, moneyWonOrLost, *dealerCards, dealerHandValue,
		drawGame, playerWins, playerHasBlackjackHand, dealerHasBlackjackHand)
}

//noinspection GoUnusedParameter
func playerGameHitCallback(deck *ds.Deck52, gameRoundData *ds.GameRoundData) {
	event.AddEvent(&event.Event{EventType: event.PlayerGameHitEvent})
}

//noinspection GoUnusedParameter
func playerGameDoubleCallback(deck *ds.Deck52, gameRoundData *ds.GameRoundData) {
	gameRoundData.GameState.Double()
}

func playerGameSplitCallback(deck *ds.Deck52, gameRoundData *ds.GameRoundData) {
	player := &gameRoundData.Player
	player.Split = true
	player.LeftHand = append(player.LeftHand, player.RightHand[1])
	player.RightHand = player.RightHand[:1]

	player.Money -= player.BetRightHand
	player.BetLeftHand = player.BetRightHand

	ui.PrintPlayerSplitInfo(ui.CardsWithValue{Cards: player.RightHand, Value: commons.GetHandValue(player.RightHand)},
		ui.CardsWithValue{Cards: player.LeftHand, Value: commons.GetHandValue(player.LeftHand)})

	gameRoundData.GameState = &PlayerWithDoubleHandState{Deck: deck, GameRoundData: gameRoundData}
	event.AddEvent(&event.Event{EventType: event.PlayerGameHitEvent})
}

//noinspection GoUnusedParameter
func playerGameStandCallback(deck *ds.Deck52, gameRoundData *ds.GameRoundData) {
	event.AddEvent(&event.Event{EventType: event.PlayerGameStandEvent})
}
