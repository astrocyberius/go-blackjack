package datastructs

const (
	RightHandActive = iota
	LeftHandActive
)

type ActiveHand uint8

// Player player data structure
type Player struct {
	Name  string
	Money TMoney

	RightHand              []PlayingCard
	RightHandPlayingDouble bool
	BetRightHand           TMoney

	LeftHand              []PlayingCard
	LeftHandPlayingDouble bool
	BetLeftHand           TMoney

	ActiveHand ActiveHand
	Split      bool
}

func (player *Player) IsRightHandActive() bool {
	return player.ActiveHand == RightHandActive
}

func (player *Player) IsLeftHandActive() bool {
	return player.ActiveHand == LeftHandActive
}

func (player *Player) MakeRightHandActive() {
	player.ActiveHand = RightHandActive
}

func (player *Player) MakeLeftHandActive() {
	player.ActiveHand = LeftHandActive
}
