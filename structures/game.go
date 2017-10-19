package structures

import (
	"time"

	"github.com/bmatsuo/hexgrid"
)

// Game object representing finite game state
type Game struct {
	ID         string // unique
	Owner      string // Account.ID of owning user -> nil == NPC ship
	LastTurn   time.Time
	Ships      []Ship
	Map        *hexgrid.Grid // https://github.com/bmatsuo/hexgrid OR https://github.com/pmcxs/hexgrid
	Credits    uint32
	Glory      int16
	ServerTurn bool
	Online     bool
	//Player2 uuid.UUID // Account.UUID of 2nd participating user
	//Player2IsFriend bool // false == foe (can attack Owner)
}

// NewGame with ships & map
func NewGame(gameID string, ownerID string) *Game {
	game := Game{
		ID:         gameID,
		Owner:      ownerID,
		Online:     true,
		Credits:    2000,
		Glory:      0,
		ServerTurn: false,
		Map:        hexgrid.NewGrid(11, 11, 128, nil, nil, nil),
	}

	return &game
}
