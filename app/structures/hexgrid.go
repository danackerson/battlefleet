package structures

import (
	"math"
)

// Point exported
type Point struct {
	X float64
	Y float64
}

// TotalHexagons returns total number of hexagons for given grid size
func TotalHexagons(gridSize int) int {
	return int(3*gridSize*(gridSize-1) + 1)
}

// original: https://github.com/gojuno/go.hexgrid/blob/master/hexgrid.go

// Hex describes a regular hexagon with Cube Coordinates (although the S coordinate is computed on the constructor)
// It's also easy to reference them as axial (trapezoidal coordinates):
// - R represents the vertical axis
// - Q the diagonal one
// - S can be ignored
// For additional reference on these coordinate systems: http://www.redblobgames.com/grids/hexagons/#coordinates
//           _ _
//         /     \
//    _ _ /(0,-1) \ _ _
//  /     \  -R   /     \
// /(-1,0) \ _ _ /(1,-1) \
// \  -Q   /     \       /
//  \ _ _ / (0,0) \ _ _ /
//  /     \       /     \
// /(-1,1) \ _ _ / (1,0) \
// \       /     \  +Q   /
//  \ _ _ / (0,1) \ _ _ /
//        \  +R   /
//         \ _ _ /
// Hex is now exported
type Hex struct {
	Q int64 // column
	R int64 // row
}

// FractionalHex is now commented
type FractionalHex struct {
	Q float64
	R float64
}

// Orientation is now commented
type Orientation struct {
	F          [4]float64
	B          [4]float64
	StartAngle float64
	Sinuses    [6]float64
	Cosinuses  [6]float64
}

// Grid is now commented
type Grid struct {
	Orientation Orientation
	Origin      Point
	Size        Point
}

// OrientationFlat means hex is flat end up: https://www.redblobgames.com/grids/hexagons/implementation.html#layout
var OrientationFlat = Orientation{
	F:          [4]float64{3.0 / 2.0, 0.0, math.Sqrt(3.0) / 2.0, math.Sqrt(3.0)}, // Forward 2x2 matrix
	B:          [4]float64{2.0 / 3.0, 0.0, -1.0 / 3.0, math.Sqrt(3.0) / 3.0},     // Backward 2x2 (inverse) matrix
	StartAngle: 0.0}

func init() {
	prehashAngles(&OrientationFlat)
}

func prehashAngles(orientation *Orientation) {
	for i := 0; i < 6; i++ {
		angle := 2.0 * math.Pi * (float64(i) + orientation.StartAngle) / 6.0
		orientation.Sinuses[i] = math.Sin(angle)
		orientation.Cosinuses[i] = math.Cos(angle)
	}
}

func round(val float64) int64 {
	if val < 0 {
		return int64(val - 0.5)
	}
	return int64(val + 0.5)
}

func max(a int64, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func min(a int64, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

// MakePoint exported
func MakePoint(x float64, y float64) Point {
	return Point{X: x, Y: y}
}

// GetX coord exported
func (point Point) GetX() float64 {
	return point.X
}

// GetY coord exported
func (point Point) GetY() float64 {
	return point.Y
}

// MakeHex now commented
func MakeHex(q int64, r int64) Hex {
	return Hex{Q: q, R: r}
}

// GetQ exported
func (hex Hex) GetQ() int64 {
	return hex.Q
}

// GetR exported
func (hex Hex) GetR() int64 {
	return hex.R
}

// GetS exported more
func (hex Hex) GetS() int64 {
	return -(hex.Q + hex.R)
}

// MakeFractionalHex is now commented
func MakeFractionalHex(q float64, r float64) FractionalHex {
	return FractionalHex{Q: q, R: r}
}

// GetQ returns the diagonal axis of the hex
func (fhex FractionalHex) GetQ() float64 {
	return fhex.Q
}

// GetR returns the vertical axis of the hex
func (fhex FractionalHex) GetR() float64 {
	return fhex.R
}

// GetS returns
func (fhex FractionalHex) GetS() float64 {
	return -(fhex.Q + fhex.R)
}

// ToHex is now commented
func (fhex FractionalHex) ToHex() Hex {
	q := round(fhex.GetQ())
	r := round(fhex.GetR())
	s := round(fhex.GetS())
	qDiff := math.Abs(float64(q) - fhex.GetQ())
	rDiff := math.Abs(float64(r) - fhex.GetR())
	sDiff := math.Abs(float64(s) - fhex.GetS())

	if qDiff > rDiff && qDiff > sDiff {
		q = -(r + s)
	} else if rDiff > sDiff {
		r = -(q + s)
	}

	return Hex{Q: q, R: r}
}

// MakeGrid exported more
func MakeGrid(origin Point, size Point) *Grid {
	return &Grid{Orientation: OrientationFlat, Origin: origin, Size: size}
}

// HexAt exported
func (grid *Grid) HexAt(point Point) Hex {
	x := (point.GetX() - grid.Origin.GetX()) / grid.Size.GetX()
	y := (point.GetY() - grid.Origin.GetY()) / grid.Size.GetY()
	q := grid.Orientation.B[0]*x + grid.Orientation.B[1]*y
	r := grid.Orientation.B[2]*x + grid.Orientation.B[3]*y
	return MakeFractionalHex(q, r).ToHex()
}

// HexCenter is creamy and smooth
func (grid *Grid) HexCenter(hex Hex) Point {
	x := (grid.Orientation.F[0]*float64(hex.GetQ())+grid.Orientation.F[1]*float64(hex.GetR()))*grid.Size.GetX() + grid.Origin.GetX()
	y := (grid.Orientation.F[2]*float64(hex.GetQ())+grid.Orientation.F[3]*float64(hex.GetR()))*grid.Size.GetY() + grid.Origin.GetY()
	return MakePoint(x, y)
}

// HexCorners are crunchy and round
func (grid *Grid) HexCorners(hex Hex) [6]Point {
	var corners [6]Point
	center := grid.HexCenter(hex)
	for i := 0; i < 6; i++ {
		x := grid.Size.GetX()*grid.Orientation.Cosinuses[i] + center.GetX()
		y := grid.Size.GetY()*grid.Orientation.Sinuses[i] + center.GetY()
		corners[i] = MakePoint(x, y)
	}
	return corners
}

// HexNeighbors lets us see all the neighbor hexes. Layers are "rings" around pt
func (grid *Grid) HexNeighbors(hex Hex, layers int64) []Hex {
	total := (layers + 1) * layers * 3
	neighbors := make([]Hex, total)
	i := 0
	for q := -layers; q <= layers; q++ {
		r1 := max(-layers, -q-layers)
		r2 := min(layers, -q+layers)
		for r := r1; r <= r2; r++ {
			if q == 0 && r == 0 {
				continue
			}

			neighbors[i] = MakeHex(q+hex.GetQ(), r+hex.GetR())
			i++
		}
	}

	// be careful and DON'T return neighbors that DON'T exist in the grid
	for i := len(neighbors) - 1; i >= 0; i-- {
		Q := neighbors[i].GetQ()
		R := neighbors[i].GetR()
		if Q < 0 {
			Q = -Q
		}
		if R < 0 {
			R = -R
		}
		if Q > int64(grid.Size.GetX()) || R > int64(grid.Size.GetY()) {
			neighbors = append(neighbors[:i], neighbors[i+1:]...)
		}
	}

	return neighbors
}

/* from https://github.com/kellydunn/golang-geo */
func intersectsWithRaycast(point Point, start Point, end Point) bool {
	if start.GetY() > end.GetY() {
		start, end = end, start
	}

	for point.GetY() == start.GetY() || point.GetY() == end.GetY() {
		newY := math.Nextafter(point.GetY(), math.Inf(1))
		point = MakePoint(point.GetX(), newY)
	}

	if point.GetY() < start.GetY() || point.GetY() > end.GetY() {
		return false
	}

	if start.GetX() > end.GetX() {
		if point.GetX() > start.GetX() {
			return false
		}
		if point.GetX() < end.GetX() {
			return true
		}
	} else {
		if point.GetX() > end.GetX() {
			return false
		}
		if point.GetX() < start.GetX() {
			return true
		}
	}

	raySlope := (point.GetY() - start.GetY()) / (point.GetX() - start.GetX())
	diagSlope := (end.GetY() - start.GetY()) / (end.GetX() - start.GetX())

	return raySlope >= diagSlope
}

func pointInGeometry(geometry [][]Point, point Point) bool {
	contains := false
	for g := 0; g < len(geometry); g++ {
		if intersectsWithRaycast(point, geometry[g][len(geometry[g])-1], geometry[g][0]) {
			contains = !contains
		}
		for i := 1; i < len(geometry[g]); i++ {
			if intersectsWithRaycast(point, geometry[g][i-1], geometry[g][i]) {
				contains = !contains
			}
		}
	}
	return contains
}
