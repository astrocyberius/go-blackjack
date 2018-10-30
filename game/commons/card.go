package commons

import (
	ds "github.com/astrocyberius/go-blackjack/datastructs"
	"math/rand"
	"time"
)

var cardSuites = [4]ds.Suit{ds.Spades, ds.Hearts, ds.Diamonds, ds.Clubs}

var cardRanks = [13]ds.Rank{ds.Ace, ds.King, ds.Queen, ds.Jack, ds.Ten, ds.Nine, ds.Eight, ds.Seven, ds.Six,
	ds.Five, ds.Four, ds.Three, ds.Two}

var cardValues = map[ds.Rank]ds.CardValue{
	ds.Ace:   11,
	ds.King:  10,
	ds.Queen: 10,
	ds.Jack:  10,
	ds.Ten:   10,
	ds.Nine:  9,
	ds.Eight: 8,
	ds.Seven: 7,
	ds.Six:   6,
	ds.Five:  5,
	ds.Four:  4,
	ds.Three: 3,
	ds.Two:   2,
}

func GetHandValue(playingCards []ds.PlayingCard) ds.CardValue {
	var value ds.CardValue = 0

	var aceCards []ds.PlayingCard
	for _, playingCard := range playingCards {
		if playingCard.Rank == ds.Ace {
			aceCards = append(aceCards, playingCard)
			continue
		}
		value += cardValues[playingCard.Rank]
	}

	for _, aceCard := range aceCards {
		if value+cardValues[aceCard.Rank] > BlackjackWinPoints {
			value += 1
		} else {
			value += cardValues[aceCard.Rank]
		}
	}

	return value
}

func IsBlackjackHand(playingCards []ds.PlayingCard, handValue ds.CardValue) bool {
	return handValue == BlackjackWinPoints && len(playingCards) == 2
}

func CreateCardDeck() *ds.Deck52 {
	deck := ds.Deck52{}

	i := 0
	for cardSuite := range cardSuites {
		for cardRank := range cardRanks {
			deck.Cards[i] = ds.PlayingCard{Rank: ds.Rank(cardRank), Suit: ds.Suit(cardSuite)}
			i++
		}
	}

	return &deck
}

func PopCardFromDeck(deck *ds.Deck52, playingCards *[]ds.PlayingCard) ds.PlayingCard {
	if len(*playingCards) == 0 {
		*playingCards = GetPlayingCards(ShuffleCardDeck(deck))
	}

	playingCard := (*playingCards)[0]
	*playingCards = append([]ds.PlayingCard{}, (*playingCards)[1:]...)
	return playingCard
}

func ShuffleCardDeck(deck *ds.Deck52) *ds.Deck52 {
	var shuffledCards [52]ds.PlayingCard

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	cards := append([]ds.PlayingCard{}, deck.Cards[:]...)
	for i := 0; i < len(shuffledCards); i++ {
		idx := r1.Intn(len(cards))
		shuffledCards[i] = cards[idx]
		cards = append(cards[0:idx], cards[idx+1:]...)
	}

	return &ds.Deck52{Cards: shuffledCards}
}

func GetPlayingCards(deck *ds.Deck52) []ds.PlayingCard {
	return append([]ds.PlayingCard{}, deck.Cards[:]...)
}
