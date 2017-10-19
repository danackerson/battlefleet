package structures

import (
	"github.com/bmatsuo/hexgrid/point"
	uuid "gopkg.in/myesui/uuid.v1"
)

// Ship object representing finite state of ship
type Ship struct {
	ID         uuid.UUID
	Owner      uuid.UUID // Account.ID or Player2.ID => nil for Server ship (NPC)
	Name       string
	Position   point.Point
	Crystals   uint32 // 0 means no firing/movement
	GunPower   uint32 // can be upgraded at Bases
	HullDamage int8   // 100% means no movement
	GunDamage  int8   // 100% means no firing
	Docked     bool   // at a base (no attack/yes repairs)?
	Type       string // [patrol,fighter,cruiser,destroyer]
	Class      string // [derelict,low,mid,high,luxury]
}
