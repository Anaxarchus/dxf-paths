package paths

import (
	"fmt"
	"math"
	"slices"

	"github.com/Anaxarchus/gdscript-libs/pkg/mathgd"
)

type Path struct {
	Closed bool
	Points []mathgd.Vector2
}

type Direction int

const (
	DirectionBackward = iota
	DirectionForward
)

func NewPathFromBulge(points []VectorB, closed bool, maxInterval float64) *Path {
	// Check if the first and the last points are the same
	if points[0].IsEqualApprox(points[len(points)-1].Vector2) {
		// Remove the last point by slicing up to the second-to-last element
		points = points[:len(points)-1]
	}

	//pts := LoopA(points, closed, maxInterval, minSteps)
	pts := []mathgd.Vector2{}
	for i := 0; i < len(points); i++ {
		nxt := mathgd.Wrapi(i+1, 0, len(points)-1)
		if points[i].B > 0 {
			arc := points[i].GetArcTo(points[nxt])
			fmt.Println("Arc: ", arc)
			pts = append(pts, arc.Discretize(maxInterval)...)
		} else {
			pts = append(pts, points[i].Vector2)
		}
	}

	return &Path{
		Points: pts,
		Closed: closed,
	}
}

func NewPath(points []mathgd.Vector2, closed bool) *Path {
	// Check if the first and the last points are the same
	if points[0] == points[len(points)-1] {
		// Remove the last point by slicing up to the second-to-last element
		points = points[:len(points)-1]
	}

	return &Path{
		Points: points,
		Closed: closed,
	}
}

func (p *Path) GetArcSnippets() []*ArcSnippet {
	arcs := []*ArcSnippet{}
	pLen := len(p.Points)
	if pLen < MinArcPoints {
		return arcs
	}

	var lastAngle float64
	var pointCount int
	var firstIdx int
	var lastIdx int

	for i := 0; i < pLen; i++ {
		currentIdx := i
		prevIdx := mathgd.Wrapi(i-1, 0, pLen)
		nextIdx := mathgd.Wrapi(i+1, 0, pLen)

		currentAngle := getAngle(p.Points[prevIdx], p.Points[currentIdx], p.Points[nextIdx])

		if math.Abs(currentAngle-lastAngle) > 0.0001 {
			if pointCount >= MinArcPoints-2 {
				arcs = append(arcs, NewArcSnippet(firstIdx, lastIdx, p))
			}
			firstIdx = currentIdx
			pointCount = 0
		}

		if currentAngle > MaxArcAngle {
			if pointCount >= MinArcPoints-2 {
				arcs = append(arcs, NewArcSnippet(firstIdx, lastIdx, p))
			}
			firstIdx = currentIdx
			pointCount = 0
		}

		if math.Abs(currentAngle) < MinArcAngle {
			if pointCount >= MinArcPoints-2 {
				arcs = append(arcs, NewArcSnippet(firstIdx, lastIdx, p))
			}
			firstIdx = currentIdx
			pointCount = 0
		}

		pointCount++

		lastIdx = currentIdx
		lastAngle = currentAngle
	}

	// Check for the last segment
	if pointCount >= MinArcPoints {
		arcs = append(arcs, NewArcSnippet(firstIdx, lastIdx, p))
	}

	return arcs
}

func getAngle(a, b, c mathgd.Vector2) float64 {
	ba := b.DirectionTo(a)
	bc := b.DirectionTo(c)
	return ba.Angle() - bc.Angle()
}

func (p *Path) AddArc(arc *Arc, maxInterval float64) {
	points := arc.Discretize(maxInterval)
	p.Points = append(p.Points, points...)
}

func (p *Path) InsertArc(arc *Arc, maxInterval float64, index int) {
	points := arc.Discretize(maxInterval)
	p.Points = slices.Insert(p.Points, index, points...)
}

func (p *Path) AddPoint(point mathgd.Vector2) {
	p.Points = append(p.Points, point)
}

func (p *Path) InsertPoint(point mathgd.Vector2, index int) {
	index = mathgd.Clampi(index, 0, len(p.Points)-1)
	p.Points = slices.Insert(p.Points, index, point)
}

func (p *Path) GetPoint(index int) mathgd.Vector2 {
	if p.Closed {
		index = mathgd.Wrapi(index, 0, len(p.Points)-1)
	}
	return p.Points[mathgd.Clampi(index, 0, len(p.Points)-1)]
}

func (p *Path) GetSegment(aIndex, bIndex int) [2]mathgd.Vector2 {
	return [2]mathgd.Vector2{p.GetPoint(aIndex), p.GetPoint(bIndex)}
}

func (p *Path) GetNearestPoint(fromPosition mathgd.Vector2) (mathgd.Vector2, int, int) {
	a, b := p.GetNearestSegment(fromPosition)
	return mathgd.GetClosestPointToSegment(fromPosition, p.GetSegment(a, b)), a, b
}

func (p *Path) GetNearestSegment(fromPosition mathgd.Vector2) (int, int) {
	var nSeg [2]int
	nDist := math.Inf(1)
	for i := 1; i < len(p.Points); i++ {
		seg := p.GetSegment(i, i-1)
		dist := mathgd.GetDistanceSquaredToSegment(fromPosition, seg)
		if dist < nDist {
			nSeg = [2]int{i - 1, i}
			nDist = dist
		}
	}
	return nSeg[0], nSeg[1]
}

func (p *Path) Offset(delta float64) error {
	mathgd.OffsetPolygon(p.Points, delta, mathgd.JoinTypeMiter)
	return nil
}

func (p *Path) Walk(start mathgd.Vector2, distance float64, direction Direction) *Snippet {
	indices, end := p.walk(start, distance, direction)
	return &Snippet{
		Indices: indices,
		Start:   start,
		End:     end,
	}
}

func (p *Path) walk(fromPosition mathgd.Vector2, distance float64, direction Direction) ([]int, mathgd.Vector2) {
	if len(p.Points) < 2 {
		return nil, fromPosition
	}

	currentPoint, a, b := p.GetNearestPoint(fromPosition)
	nextIndex := a
	indices := []int{}
	if currentPoint.IsEqualApprox(p.Points[a]) {
		indices = append(indices, a)
	} else if currentPoint.IsEqualApprox(p.Points[b]) {
		indices = append(indices, b)
		nextIndex = b
	}
	distanceRemaining := distance

	for distanceRemaining > 0 {
		nextIndex = p.Next(nextIndex, direction)

		nextPoint := p.Points[nextIndex]
		distanceToNext := currentPoint.DistanceTo(nextPoint)

		if distanceToNext > distanceRemaining {
			// Move partially to the next point
			t := distanceRemaining / distanceToNext
			endPosition := mathgd.Vector2{
				X: currentPoint.X + t*(nextPoint.X-currentPoint.X),
				Y: currentPoint.Y + t*(nextPoint.Y-currentPoint.Y),
			}
			return indices, endPosition
		}

		// Move fully to the next point
		currentPoint = nextPoint
		indices = append(indices, nextIndex)
		distanceRemaining -= distanceToNext

		if !p.Closed && (nextIndex == 0 || nextIndex == len(p.Points)-1) {
			break
		}
	}

	return indices, currentPoint
}

func (p *Path) Next(index int, direction Direction) int {
	if direction == DirectionForward {
		index += 1
	} else {
		index -= 1
	}
	if p.Closed {
		return mathgd.Wrapi(index, 0, len(p.Points)-1)
	} else {
		return mathgd.Clampi(index, 0, len(p.Points)-1)
	}
}

func (p *Path) Previous(index int, direction Direction) int {
	if direction == DirectionForward {
		index -= 1
	} else {
		index += 1
	}
	if p.Closed {
		return mathgd.Wrapi(index, 0, len(p.Points)-1)
	} else {
		return mathgd.Clampi(index, 0, len(p.Points)-1)
	}
}
