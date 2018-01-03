package structures

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// NewGameUUID is the default value for creating a new game
var NewGameUUID = "__new__"

// GridSize is # of "rings" radiating out from Origin hex (size of "5" => 91 hexagons in total)
// This setting is also used to control Client side game grid rendering (game.tmpl => vg.Scene)
var GridSize = float64(5)

// Game object representing finite game state
type Game struct {
	ID         string        // unique
	Owner      bson.ObjectId `bson:"_id,omitempty"` // Account.ID of owning user -> nil == NPC ship
	LastTurn   time.Time
	Ships      []*Ship
	Map        *Grid
	Credits    uint32
	Glory      int16
	ServerTurn bool
	Online     bool
	//Player2 uuid.UUID // Account.UUID of 2nd participating user
	//Player2IsFriend bool // false == foe (can attack Owner)
}

// NewGame with ships & map
func NewGame(gameID string, ownerID bson.ObjectId) *Game {
	ship := &Ship{
		ID:         gameID,
		Owner:      ownerID,
		Name:       "StarWars",
		Position:   MakePoint(0, 0), // Origin (center hex) on Grid
		Crystals:   100,
		GunPower:   10,
		HullDamage: 0,
		GunDamage:  0,
		Docked:     true,
		Type:       "patrol",
		Class:      "low",
	}

	origin := MakePoint(0, 0)
	size := MakePoint(GridSize, GridSize)
	grid := MakeGrid(origin, size)

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
