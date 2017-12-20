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

// TestNeighbors now commented
func TestNeighbors(t *testing.T) {
	/*grid := MakeGrid(OrientationFlat, MakePoint(10, 20), MakePoint(20, 10))
	hex := grid.HexAt(MakePoint(666, 666))
	expectedNeighbors := [18]int64{
		920, 922, 944, 915, 921, 923, 945, 916, 918,
		926, 948, 917, 919, 925, 927, 960, 962, 968}
	neighbors := grid.HexNeighbors(hex, 2)*/
	// TODO: check it
}
