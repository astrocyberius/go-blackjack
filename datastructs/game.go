package datastructs

// TMoney money type
type TMoney float64

// GameRoundData game round data
type GameRoundData struct {
	PlayingCards []PlayingCard
	DealerCards  []PlayingCard
	Player       Player
	GameState    GameState
}

// GameState the game state machine
type GameState interface {
	Hit()
	Double()
	Stand()
}
