package paths

import (
	"math"

	"github.com/Anaxarchus/gdscript-libs/pkg/mathgd"
)

type BulgeVector interface {
	GetX() float64
	GetY() float64
	GetBulge() float64
}

type VectorB struct {
	mathgd.Vector2
	B float64
}

func (vb VectorB) GetArcTo(to VectorB) *Arc {
	return NewArc(bulgeToArc(to.Vector2, vb.Vector2, vb.B))
}

func bulgeToArc(p1, p2 mathgd.Vector2, bulge float64) (mathgd.Vector2, float64, float64, float64) {
	// Calculate the distance between the points
	distance := p1.DistanceTo(p2)

	// Calculate the radius of the arc
	b := bulge
	r := (distance * (1 + b*b)) / (4 * b)

	// Calculate the angle from point 1 to point 2
	angleP1P2 := math.Atan2(p2.Y-p1.Y, p2.X-p1.X)
	//angleP1P2 := p2.AngleTo(p1)

	// Calculate the center of the arc
	theta := angleP1P2 - (math.Pi/2 - 2*math.Atan(b))
	cX := p1.X + (r * math.Cos(theta))
	cY := p1.Y + (r * math.Sin(theta))
	center := mathgd.NewVector2(cX, cY)

	// Calculate the start and end angles
	var startAngle, endAngle float64
	if b > 0 {
		startAngle = math.Atan2(p2.Y-cY, p2.X-cX)
		endAngle = math.Atan2(p1.Y-cY, p1.X-cX)
	} else {
		startAngle = math.Atan2(p1.Y-cY, p1.X-cX)
		endAngle = math.Atan2(p2.Y-cY, p2.X-cX)
	}

	return center, math.Abs(r), startAngle, endAngle
}
