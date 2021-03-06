package structures

import (
	"gopkg.in/mgo.v2/bson"
)

// Ship object representing finite state of ship
type Ship struct {
	ID         string
	Owner      bson.ObjectId `bson:"_id,omitempty"` // Account.ID or Player2.ID => nil for Server ship (NPC)
	Name       string
	Position   Point
	Crystals   uint32 // 0 means no firing/movement
	GunPower   uint32 // can be upgraded at Bases
	HullDamage int8   // 100% means no movement
	GunDamage  int8   // 100% means no firing
	Docked     bool   // at a base (no attack/yes repairs)?
	Type       string // [patrol,fighter,cruiser,destroyer]
	Class      string // [derelict,low,mid,high,luxury]
}
