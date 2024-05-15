package paths

import (
	"math"

	"github.com/Anaxarchus/gdscript-libs/pkg/mathgd"
)

type Arc struct {
	Origin     mathgd.Vector2
	Radius     float64
	StartAngle float64
	EndAngle   float64
}

const MaxArcAngle float64 = 22.5
const MinArcAngle float64 = mathgd.CMP_EPSILON2
const MinArcPoints int = 17
const MinArcSegments int = 16

func NewArc(origin mathgd.Vector2, radius, start, end float64) *Arc {
	return &Arc{
		Origin:     origin,
		Radius:     radius,
		StartAngle: start,
		EndAngle:   end,
	}
}

func NewArcFromArcPoints(points []mathgd.Vector2) *Arc {
	origin, radius := fitArc(points)
	a1 := points[0]
	a2 := points[len(points)-1]
	return &Arc{
		Origin:     origin,
		Radius:     radius,
		StartAngle: origin.AngleToPoint(a1),
		EndAngle:   origin.AngleToPoint(a2),
	}
}

func (a Arc) GetArcBetweenPoints(start, end mathgd.Vector2) *Arc {
	a.StartAngle = a.AngleToPoint(start)
	a.EndAngle = a.AngleToPoint(end)
	return &a
}

func (a *Arc) Project(point mathgd.Vector2) mathgd.Vector2 {
	return a.Origin.Add(a.Origin.DirectionTo(point).Mulf(a.Radius))
}

func (a *Arc) AngleToPoint(point mathgd.Vector2) float64 {
	return a.Origin.AngleToPoint(point)
}

func (a *Arc) Discretize(maxInterval float64) []mathgd.Vector2 {
	// Calculate the total angle span of the arc
	totalAngle := a.EndAngle - a.StartAngle
	if totalAngle < 0 {
		totalAngle += 2 * math.Pi
	}

	// Calculate the total length of the arc
	arcLength := a.Radius * totalAngle

	// Calculate the minimum number of steps required based on maxInterval
	minStepsBasedOnInterval := int(math.Ceil(arcLength / maxInterval))

	// Ensure the number of steps is at least minSteps and is odd
	numSteps := minStepsBasedOnInterval
	if numSteps < MinArcPoints {
		numSteps = MinArcPoints
	}
	if numSteps%2 == 0 {
		numSteps++
	}

	// Calculate the angle step
	angleStep := totalAngle / float64(numSteps-1)

	// Generate the points along the arc
	points := make([]mathgd.Vector2, 0, numSteps)
	for i := 0; i < numSteps; i++ {
		angle := a.StartAngle + float64(i)*angleStep
		x := a.Origin.X + a.Radius*math.Cos(angle)
		y := a.Origin.Y + a.Radius*math.Sin(angle)
		points = append(points, mathgd.Vector2{X: x, Y: y})
	}

	return points
}

// ArcDirection determines the direction of the arc: 1 for clockwise, -1 for counter-clockwise.
// It takes lastPosition, currentPosition and arcOrigin, all as Vector3.
func ArcDirection(lastPosition, position, arcOrigin mathgd.Vector2) int {
	// Convert to coordinates relative to arcOrigin
	lastRel := lastPosition.Sub(arcOrigin)
	currentRel := position.Sub(arcOrigin)

	// Calculate angles from arcOrigin
	lastAngle := lastRel.Angle()
	currentAngle := currentRel.Angle()

	// Calculate angle difference and determine direction
	angleDifference := currentAngle - lastAngle

	// Normalize the angle to be within the range -π to π
	if angleDifference > math.Pi {
		angleDifference -= 2 * math.Pi
	} else if angleDifference < -math.Pi {
		angleDifference += 2 * math.Pi
	}

	// Determine the direction based on the angle difference
	if angleDifference > 0 {
		return -1 // CCW
	} else {
		return 1 // CW
	}
}

// ProcessArcs calculates the center and radius of a circle that passes through the first,
// last, and a nearest point to the midpoint of the first and last points from a given
// slice of Vector2 points that represent an arc.
func fitArc(arcPoints []mathgd.Vector2) (mathgd.Vector2, float64) {
	e1 := arcPoints[0]                // Start point of the arc
	e2 := arcPoints[len(arcPoints)-1] // End point of the arc

	// Calculate the midpoint of the line segment from e1 to e2
	mp := e1.Add(e2).Divf(2.0)

	// Find the point nearest to the midpoint among the arc points
	minDist := math.MaxFloat64
	var np mathgd.Vector2
	for _, point := range arcPoints {
		dist := mp.DistanceTo(point)
		if dist < minDist {
			minDist = dist
			np = point
		}
	}

	// Calculate the distances of the sides of the triangle formed by e1, e2, and np
	a := e1.DistanceTo(e2) // Distance between end points
	b := e1.DistanceTo(np) // Distance from start point to nearest point
	c := e2.DistanceTo(np) // Distance from end point to nearest point

	// Calculate the radius of the circumscribed circle using Heron's formula
	// and the area of a triangle calculation
	product := (a + b + c) * (a + b - c) * (a - b + c) * (b + c - a)
	if product <= 0 {
		return np, 0 // If product is non-positive, radius calculation fails
	}
	rad := (a * b * c) / math.Sqrt(product)

	// The origin of the circle is calculated by extending the vector from np to mp
	// by the calculated radius
	center := np.DirectionTo(mp).Mulf(rad)

	return center, rad
}
