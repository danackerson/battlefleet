package structures

import (
	"math"
	"testing"
)

func validateHex(t *testing.T, e Hex, r Hex) {
	if e.GetQ() != r.GetQ() || e.GetR() != r.GetR() {
		t.Errorf("expected hex{q: %d, r: %d} but got hex{q: %d, r: %d}", e.GetQ(), e.GetR(), r.GetQ(), r.GetR())
	}
}

// TestFlat now commented
func TestFlat(t *testing.T) {
	grid := MakeGrid(MakePoint(10, 20), MakePoint(20, 10))
	validateHex(t, MakeHex(0, 37), grid.HexAt(MakePoint(13, 666)))
	validateHex(t, MakeHex(22, -11), grid.HexAt(MakePoint(666, 13)))
	validateHex(t, MakeHex(-1, -39), grid.HexAt(MakePoint(-13, -666)))
	validateHex(t, MakeHex(-22, 9), grid.HexAt(MakePoint(-666, -13)))
}

func validatePoint(t *testing.T, e Point, r Point, precision float64) {
	if math.Abs(e.GetX()-r.GetX()) > precision || math.Abs(e.GetY()-r.GetY()) > precision {
		t.Errorf("expected point{x: %f, y: %f} but got point{x: %f, y: %f}", e.GetX(), e.GetY(), r.GetX(), r.GetY())
	}
}

func TestCoordinatesFlat(t *testing.T) {
	grid := MakeGrid(MakePoint(10, 20), MakePoint(20, 10))
	hex := grid.HexAt(MakePoint(666, 666))
	validatePoint(t, MakePoint(670.00000, 660.85880), grid.HexCenter(hex), 0.00001)
	expectedCorners := [6]Point{
		MakePoint(690.00000, 660.85880),
		MakePoint(680.00000, 669.51905),
		MakePoint(660.00000, 669.51905),
		MakePoint(650.00000, 660.85880),
		MakePoint(660.00000, 652.19854),
		MakePoint(680.00000, 652.19854)}
	corners := grid.HexCorners(hex)
	for i := 0; i < 6; i++ {
		validatePoint(t, expectedCorners[i], corners[i], 0.00001)
	}
}

// TestNeighborsOfOuterEdge verifies all Hex cells around outer left tile {-5,5}
func TestNeighborsOfOuterEdge(t *testing.T) {
	// 5 cells/tiles radiating outward from Origin point (0,0)
	// make 5 rings total and 1+6+12+18+24+30=91 cells/tiles
	size := 5
	origin := MakePoint(0, 0)

	grid := MakeGrid(origin, MakePoint(float64(size), float64(size)))
	hex := grid.HexAt(MakePoint(-35, 25)) // outer, upper left corner point on grid
	neighbors := grid.HexNeighbors(hex, 1)

	// around outside corner tile, we have 3 tiles around it
	expectedNeighbors := 3
	if expectedNeighbors != len(neighbors) {
		t.Errorf("expected Grid of size %d to have %d neighbors but got %d", size, expectedNeighbors, len(neighbors))
	}

	firstHex := &Hex{Q: -5, R: 4}
	lastHex := &Hex{Q: -4, R: 5}
	if neighbors[0] != *firstHex {
		t.Errorf("expected first neighbor tile to be Hex{Q:-5, R:4} but got Hex{Q:%d, R:%d}\n", neighbors[0].GetQ(), neighbors[0].GetR())
	}
	if neighbors[expectedNeighbors-1] != *lastHex {
		t.Errorf("expected last neighbor tile to be Hex{Q:-4, R:5} but got Hex{Q:%d, R:%d}\n", neighbors[0].GetQ(), neighbors[0].GetR())
	}

	/* All tiles in size 5 grid
	for index, hex := range neighbors {
		log.Printf("neighbor[%d]: Hex {Q: %d, R: %d}\n", index, hex.GetQ(), hex.GetR())
	}*/
}

// TestNeighborsOfOrigin verifies all Hex cells around game grid
func TestNeighborsOfOrigin(t *testing.T) {
	// 5 cells/tiles radiating outward from Origin point (0,0)
	// make 5 rings total and 1+6+12+18+24+30=91 cells/tiles
	size := 5
	origin := MakePoint(0, 0)

	grid := MakeGrid(origin, MakePoint(float64(size), float64(size)))
	hex := grid.HexAt(origin) // Origin point Q:0,R:0 in middle of grid
	neighbors := grid.HexNeighbors(hex, int64(size))

	// in Grid of size 5, we have 90 tiles around Origin tile
	expectedNeighbors := 90
	if expectedNeighbors != len(neighbors) {
		t.Errorf("expected Grid of size %d to have %d neighbors but got %d", size, expectedNeighbors, len(neighbors))
	}
	firstHex := &Hex{Q: -5, R: 0}
	lastHex := &Hex{Q: 5, R: 0}
	if neighbors[0] != *firstHex {
		t.Errorf("expected first neighbor tile to be Hex{Q:-5, R:0} but got Hex{Q:%d, R:%d}\n", neighbors[0].GetQ(), neighbors[0].GetR())
	}
	if neighbors[expectedNeighbors-1] != *lastHex {
		t.Errorf("expected last neighbor tile to be Hex{Q:5, R:0} but got Hex{Q:%d, R:%d}\n", neighbors[0].GetQ(), neighbors[0].GetR())
	}

	/* All tiles in size 5 grid
	for index, hex := range neighbors {
		log.Printf("neighbor[%d]: Hex {Q: %d, R: %d}\n", index, hex.GetQ(), hex.GetR())
	}*/
}
