package structures

import (
	"time"

	"github.com/danackerson/battlefleet/hexgrid"
)

// Game object representing finite game state
type Game struct {
	ID         string // unique
	Owner      string // Account.ID of owning user -> nil == NPC ship
	LastTurn   time.Time
	Ships      []*Ship
	Map        *hexgrid.Grid
	Credits    uint32
	Glory      int16
	ServerTurn bool
	Online     bool
	//Player2 uuid.UUID // Account.UUID of 2nd participating user
	//Player2IsFriend bool // false == foe (can attack Owner)
}

// NewGame with ships & map
func NewGame(gameID string, ownerID string) *Game {
	ship := &Ship{
		ID:         gameID,
		Owner:      ownerID,
		Name:       "StarWars",
		Position:   hexgrid.MakePoint(0, 0),
		Crystals:   100,
		GunPower:   10,
		HullDamage: 0,
		GunDamage:  0,
		Docked:     true,
		Type:       "patrol",
		Class:      "low",
	}

	center := hexgrid.MakePoint(0, 0)
	size := hexgrid.MakePoint(11, 11)
	grid := hexgrid.MakeGrid(hexgrid.OrientationFlat, center, size)

	game := Game{
		ID:         gameID,
		Owner:      ownerID,
		Online:     true,
		Credits:    2000,
		Glory:      0,
		ServerTurn: false,
		Map:        grid,
		Ships:      []*Ship{ship},
		LastTurn:   time.Now(),
	}

	//spew.Dump(game.Map)
	return &game
}
