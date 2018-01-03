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
	size := GridSize
	origin := MakePoint(0, 0)

	grid := MakeGrid(origin, MakePoint(float64(size), float64(size)))
	hex := &Hex{-int64(GridSize), int64(GridSize)} // outer, upper left corner point on grid
	neighbors := grid.HexNeighbors(*hex, 1)

	// around outside corner tile, we have 3 tiles around it
	expectedNeighbors := 3
	if expectedNeighbors != len(neighbors) {
		t.Errorf("expected Grid of size %f to have %d neighbors but got %d", size, expectedNeighbors, len(neighbors))
	}

	firstHex := &Hex{Q: -int64(GridSize), R: int64(GridSize) - 1}
	lastHex := &Hex{Q: -int64(GridSize - 1), R: int64(GridSize)}
	if neighbors[0] != *firstHex {
		t.Errorf("expected first neighbor tile to be Hex{Q:-%f, R:%f} but got Hex{Q:%d, R:%d}\n", GridSize, (GridSize - 1), neighbors[0].GetQ(), neighbors[0].GetR())
	}
	if neighbors[expectedNeighbors-1] != *lastHex {
		t.Errorf("expected last neighbor tile to be Hex{Q:-%f, R:%f} but got Hex{Q:%d, R:%d}\n", (GridSize - 1), GridSize, neighbors[0].GetQ(), neighbors[0].GetR())
	}

	/* All tiles in size 5 grid
	for index, hex := range neighbors {
		log.Printf("neighbor[%d]: Hex {Q: %d, R: %d}\n", index, hex.GetQ(), hex.GetR())
	}*/
}

// TestNeighborsOfOrigin verifies all Hex cells around game grid
func TestNeighborsOfOrigin(t *testing.T) {
	// 5 cells/tiles radiating outward from Origin point (0,0)
	// make 5 rings total and 1+6+12+18+24+30 +36+42+48+54+60+66=91 cells/tiles
	origin := MakePoint(0, 0)

	grid := MakeGrid(origin, MakePoint(float64(GridSize), float64(GridSize)))
	hex := grid.HexAt(origin) // Origin point Q:0,R:0 in middle of grid
	neighbors := grid.HexNeighbors(hex, int64(GridSize))

	// in Grid of size 5, we have 90 tiles around Origin tile
	expectedNeighbors := TotalHexagons(int(GridSize)+1) - 1
	if expectedNeighbors != len(neighbors) {
		t.Errorf("expected Grid of size %f to have %d neighbors but got %d", GridSize, expectedNeighbors, len(neighbors))
	}
	firstHex := &Hex{Q: -int64(GridSize), R: 0}
	lastHex := &Hex{Q: int64(GridSize), R: 0}
	if neighbors[0] != *firstHex {
		t.Errorf("expected first neighbor tile to be Hex{Q:-%f, R:0} but got Hex{Q:%d, R:%d}\n", GridSize, neighbors[0].GetQ(), neighbors[0].GetR())
	}
	if neighbors[expectedNeighbors-1] != *lastHex {
		t.Errorf("expected last neighbor tile to be Hex{Q:%f, R:0} but got Hex{Q:%d, R:%d}\n", GridSize, neighbors[0].GetQ(), neighbors[0].GetR())
	}

	/* All tiles in size 5 grid
	for index, hex := range neighbors {
		log.Printf("neighbor[%d]: Hex {Q: %d, R: %d}\n", index, hex.GetQ(), hex.GetR())
	}*/
}
