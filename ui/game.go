package ui

import (
	"bufio"
	"fmt"
	ds "github.com/astrocyberius/go-blackjack/datastructs"
	"os"
	"time"
)

// PrintIntro prints the introduction text
func PrintIntro() {
	ClearScreen()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                 Welcome to Blackjack Game                    ")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()
	fmt.Println()
}

// GetPlayerName reads the player's name from console and returns it
func GetPlayerName() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your name: ")
	name := ""
	if scanner.Scan() {
		name = scanner.Text()
	}
	return name
}

// PrintWelcomePlayer prints the welcome player text
func PrintWelcomePlayer(name string) {
	ClearScreen()
	fmt.Printf("Welcome %v.\n", name)
	printPressEnterToContinue()
}

// MenuOptionCallback menu callback methods
type MenuOptionCallback struct {
	Play func()
	Info func()
	Quit func()
}

// PrintMenuAndHandlePlayerInput prints the menu and handles the chosen option
func PrintMenuAndHandlePlayerInput(menuOptionCallBack *MenuOptionCallback) {
	ClearScreen()

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                              Menu                            ")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                       1. Play the game                       ")
	fmt.Println("                       2. Show game info                      ")
	fmt.Println("                       3. Quit the game                       ")

	option := readUint8("Option: ", 1, 3)

	switch option {
	case 1:
		menuOptionCallBack.Play()
	case 2:
		menuOptionCallBack.Info()
	case 3:
		menuOptionCallBack.Quit()
	}
}

// GetBetOptions function to calculate the bet options from the provided money
type GetBetOptions func(playerMoney ds.TMoney) []ds.TMoney

// GetPlayerBetInput retrieves the player's provided bet
func GetPlayerBetInput(getBetOptions GetBetOptions, playerMoney, minimumBetValue, maximumBetValue ds.TMoney) (ds.TMoney, bool) {
	ClearScreen()

	var bet ds.TMoney = 0
	deal := false

	for deal != true {
		fmt.Println("--------------------------------------------------------------")
		fmt.Println("                      Place Your Bets                         ")
		fmt.Println("--------------------------------------------------------------")

		fmt.Printf("Cash: €%.2f Bet: %v\n", playerMoney, bet)
		fmt.Println()
		fmt.Println("Select from the following options:")
		fmt.Println("1. Place bet")
		fmt.Println("2. Deal")
		fmt.Println("3. Clear bet")
		fmt.Println("4. Quit to main menu")

		option := readUint8("Option: ", 1, 4)
		fmt.Printf("\n\n")

		switch option {
		case 1:
			moneyAllowedToSpend := getMoneyAllowedToSpend(playerMoney, bet, maximumBetValue)

			if moneyAllowedToSpend < minimumBetValue {
				if playerMoney < minimumBetValue {
					fmt.Println("You don't have enough money.")
				} else {
					fmt.Println("You have reached the maximum allowed bet.")
				}
				printPressEnterToContinue()
			} else {
				betInput := handlePlayBetInput(getBetOptions(moneyAllowedToSpend), moneyAllowedToSpend, minimumBetValue, maximumBetValue)
				bet += betInput
				playerMoney -= betInput
			}
		case 2:
			if bet == 0 {
				fmt.Println("Place a bet first.")
				printPressEnterToContinue()
			} else {
				deal = true
			}
		case 3:
			playerMoney += bet
			bet = 0
		case 4:
			return 0, true
		}

		ClearScreen()
	}

	return bet, false
}

func getMoneyAllowedToSpend(playerMoney, bet, maximumBetValue ds.TMoney) ds.TMoney {
	moneyToSpend := playerMoney
	if moneyToSpend > maximumBetValue-bet {
		moneyToSpend = maximumBetValue - bet
	}

	return moneyToSpend
}

func handlePlayBetInput(betOptions []ds.TMoney, moneyAllowedToSpend, minimalBetValue, maximumBetValue ds.TMoney) ds.TMoney {
	var bet ds.TMoney = 0

	fmt.Printf("Minimum Bet: %v Maximum Bet: %v Allowed To Spend: %v\n", minimalBetValue, maximumBetValue, moneyAllowedToSpend)
	fmt.Println("The available bet options are:")
	for i, v := range betOptions {
		option := i + 1
		fmt.Printf("%v. %v\n", option, v)
	}

	option := readUint8("Bet option: ", 1, len(betOptions))

	chip := betOptions[option-1]
	fmt.Printf("How much do you want to bet of chip %v?\n", chip)
	option = readUint8("Number of chips: ", 1, int(moneyAllowedToSpend/chip))

	bet = chip * ds.TMoney(option)

	return bet
}

// PlayerGameOptionCallback player game callback methods
type PlayerGameOptionCallback struct {
	Hit    func(deck *ds.Deck52, gameRoundData *ds.GameRoundData)
	Stand  func(deck *ds.Deck52, gameRoundData *ds.GameRoundData)
	Double func(deck *ds.Deck52, gameRoundData *ds.GameRoundData)
	Split  func(deck *ds.Deck52, gameRoundData *ds.GameRoundData)
}

// CardsWithValue cards information with the total value
type CardsWithValue struct {
	Cards []ds.PlayingCard
	Value ds.CardValue
}

// PrintGameOptionsAndHandlePlayerInput print game options and handle player input
func PrintGameOptionsAndHandlePlayerInput(playerGameOptionCallback *PlayerGameOptionCallback,
	deck *ds.Deck52, gameRoundData *ds.GameRoundData,
	player *ds.Player, playerCardsWithValue CardsWithValue, dealerCardsWithValue CardsWithValue,
	withDoubleOption bool, withSplitOption bool) {

	ClearScreen()

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                       Round Options                          ")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println()

	playerCards := playerCardsWithValue.Cards
	playerHandValue := playerCardsWithValue.Value
	dealerCard := dealerCardsWithValue.Cards[0]
	dealerHandValue := dealerCardsWithValue.Value

	fmt.Println("      House        ")
	fmt.Println("-------------------")
	fmt.Printf("Dealer's #1 card is %v of %v.\n", dealerCard.RankLabel(), dealerCard.SuitLabel())
	fmt.Printf("Dealer's hand value is %v.\n", dealerHandValue)
	fmt.Println()

	fmt.Println("        You        ")
	fmt.Println("-------------------")

	if player.Split {
		if player.IsRightHandActive() {
			fmt.Println("You are aplying with your right hand cards.")
		} else {
			fmt.Println("You are aplying with your left hand cards.")
		}
	}
	fmt.Println()

	for i, playerCard := range playerCards {
		fmt.Printf("Your #%v card is %v of %v.\n", i+1, playerCard.RankLabel(), playerCard.SuitLabel())
	}
	fmt.Printf("Your hand value is %v.\n", playerHandValue)
	fmt.Println()

	fmt.Println("      Options      ")
	fmt.Println("-------------------")
	fmt.Println()

	callBackFunctions := []func(deck *ds.Deck52, gameRoundData *ds.GameRoundData){playerGameOptionCallback.Hit}
	fmt.Println("Select from the following options:")
	fmt.Println("1. Hit")

	maxOptionValue := 2

	if withDoubleOption {
		callBackFunctions = append(callBackFunctions, playerGameOptionCallback.Double)
		fmt.Printf("%v. Double\n", maxOptionValue)
		maxOptionValue++
	}

	if withSplitOption {
		callBackFunctions = append(callBackFunctions, playerGameOptionCallback.Split)
		fmt.Printf("%v. Split\n", maxOptionValue)
		maxOptionValue++
	}

	callBackFunctions = append(callBackFunctions, playerGameOptionCallback.Stand)
	fmt.Printf("%v. Stand\n", maxOptionValue)

	option := readUint8("Option: ", 1, maxOptionValue)

	callBackFunctions[option-1](deck, gameRoundData)
}

// PrintSwitchingToPlayerLeftHand print switching to player's left hand cards.
func PrintSwitchingToPlayerLeftHand() {
	ClearScreen()
	fmt.Println("Switching to your left hand.")
	time.Sleep(1500 * time.Millisecond)
}

// PrintPlayerSplitInfo print player split info
func PrintPlayerSplitInfo(playerRightHandWithValue CardsWithValue, playerLeftHandWithValue CardsWithValue) {
	ClearScreen()

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                    Your Right Hand Cards                     ")
	fmt.Println("--------------------------------------------------------------")
	printPlayerHand(playerRightHandWithValue)

	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                    Your Left Hand Cards                      ")
	fmt.Println("--------------------------------------------------------------")
	printPlayerHand(playerLeftHandWithValue)
	fmt.Println()

	printPressEnterToContinue()
}

func printPlayerHand(playerHandWithValue CardsWithValue) {
	for i, playerCard := range playerHandWithValue.Cards {
		fmt.Printf("Your #%v card is %v of %v.\n", i+1, playerCard.RankLabel(), playerCard.SuitLabel())
	}
	fmt.Printf("Your hand value is %v.\n", playerHandWithValue.Value)
}

// PrintPlayerHitCard print player hit card
func PrintPlayerHitCard(card ds.PlayingCard) {
	ClearScreen()
	fmt.Printf("You received card %v of %v.\n", card.RankLabel(), card.SuitLabel())
	time.Sleep(1500 * time.Millisecond)
}

// PrintDealerHitCard print dealer hit card
func PrintDealerHitCard(card ds.PlayingCard) {
	ClearScreen()
	fmt.Printf("Dealer received card %v of %v.\n", card.RankLabel(), card.SuitLabel())
	time.Sleep(1500 * time.Millisecond)
}

// PrintRoundOutcome print game outcome
func PrintRoundOutcome(playerCards []ds.PlayingCard, playerHandValue ds.CardValue, playerMoney, moneyWonOrLost ds.TMoney, dealerCards []ds.PlayingCard, dealerHandValue ds.CardValue,
	drawGame, playerWins, playerHasBlackjackHand, dealerHasBlackjackHand bool) {

	ClearScreen()

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("                       Round Outcome                          ")
	fmt.Println("--------------------------------------------------------------")

	// Simulate the dealer taking a card
	for i, dealerCard := range dealerCards {
		fmt.Printf("Dealer's #%v card is %v of %v.\n", i+1, dealerCard.RankLabel(), dealerCard.SuitLabel())
	}
	fmt.Printf("Dealer's hand value is %v.\n", dealerHandValue)
	fmt.Println()

	for i, playerCard := range playerCards {
		fmt.Printf("Your #%v card is %v of %v.\n", i+1, playerCard.RankLabel(), playerCard.SuitLabel())
	}
	fmt.Printf("Your hand value is %v.\n", playerHandValue)

	fmt.Println()
	if drawGame {
		fmt.Println("The game ends in a draw.")
	} else if playerWins {
		if playerHasBlackjackHand {
			fmt.Println("Your win with blackjack hand.")
		} else {
			fmt.Printf("You win with hand value %v.\n", playerHandValue)
		}
		fmt.Printf("You win €%.2f money.\n", moneyWonOrLost)
	} else {
		if dealerHasBlackjackHand {
			fmt.Println("House wins with blackjack hand.")
		} else {
			fmt.Printf("House wins with hand value %v.\n", dealerHandValue)
		}
		fmt.Printf("You lose €%.2f money.\n", moneyWonOrLost)
	}
	fmt.Println()
	fmt.Printf("You have €%.2f money left.\n", playerMoney)
	printPressEnterToContinue()
}

// PrintNotEnoughMoneyLeftToPlay prints not enough money left to play
func PrintNotEnoughMoneyLeftToPlay() {
	ClearScreen()
	fmt.Println("You don't have enough money left to play another round.")
	printPressEnterToContinue()
}

// PrintInfo prints the game info
func PrintInfo() {
	ClearScreen()

	fmt.Println("The blackjack info")
	printPressEnterToContinue()
	ClearScreen()
}

// PrintQuit prints the quit text
func PrintQuit() {
	ClearScreen()
	fmt.Println("I hope you enjoyed playing this blackjack game. See you next time. Bye.")
	printPressEnterToContinue()
}

// PrintMessage prints a message on the screen
func PrintMessage(message string, clearScreen bool, pressEnterToContinue bool) {
	if clearScreen {
		ClearScreen()
	}

	fmt.Println(message)

	if pressEnterToContinue {
		printPressEnterToContinue()
	}
}

func printPressEnterToContinue() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func readUint8(prompt string, minOptionValue, maxOptionValue int) uint8 {
	stdin := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	var input uint8
	for {
		_, err := fmt.Fscan(stdin, &input)
		if err == nil && int(input) >= minOptionValue && int(input) <= maxOptionValue {
			break
		}
		stdin.ReadString('\n')
		fmt.Print("Sorry, invalid input. Try again: ")
	}
	return input
}
