package datastructs

// Rank the card rank
type Rank uint8

// Suit the card suit
type Suit uint8

// CardValue the value of the card
//noinspection GoNameStartsWithPackageName
type CardValue uint8

// The ranks of the card
const (
	Ace = iota
	King
	Queen
	Jack
	Ten
	Nine
	Eight
	Seven
	Six
	Five
	Four
	Three
	Two
)

// The suits of the card
const (
	Spades = iota
	Hearts
	Diamonds
	Clubs
)

var cardRankLabel = map[Rank]string{
	Ace:   "Ace",
	King:  "King",
	Queen: "Queen",
	Jack:  "Jack",
	Ten:   "Ten",
	Nine:  "Nine",
	Eight: "Eight",
	Seven: "Seven",
	Six:   "Six",
	Five:  "Five",
	Four:  "Four",
	Three: "Three",
	Two:   "Two",
}

var cardSuitLabel = map[Suit]string{
	Spades:   "Spades",
	Hearts:   "Hearts",
	Diamonds: "Diamonds",
	Clubs:    "Clubs",
}

// PlayingCard the playing card
type PlayingCard struct {
	Rank Rank
	Suit Suit
}

// Deck52 52-card deck
type Deck52 struct {
	Cards [52]PlayingCard
}

// RankLabel returns the label of the rank
func (playingCard *PlayingCard) RankLabel() string {
	return cardRankLabel[playingCard.Rank]
}

// SuitLabel returns the label of the suit
func (playingCard *PlayingCard) SuitLabel() string {
	return cardSuitLabel[playingCard.Suit]
}
