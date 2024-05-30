package paths

import (
	"errors"
	"math"
	"slices"

	gd "github.com/Anaxarchus/zero-gdscript"
	"github.com/Anaxarchus/zero-gdscript/pkg/geometry2d"
	"github.com/Anaxarchus/zero-gdscript/pkg/vector2"
)

type Path struct {
	Closed   bool
	Reversed bool
	Points   []vector2.Vector2
}

type Direction int

const (
	DirectionBackward Direction = iota
	DirectionForward
)

func NewPath(points []vector2.Vector2, closed bool) *Path {
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

func (p *Path) AddPoint(point vector2.Vector2) {
	p.Points = append(p.Points, point)
}

func (p *Path) InsertPoint(point vector2.Vector2, index int) {
	index = gd.Clampi(index, 0, len(p.Points)-1)
	p.Points = slices.Insert(p.Points, index, point)
}

func (p *Path) GetPoint(index int) vector2.Vector2 {
	if p.Closed {
		index = gd.Wrapi(index, 0, len(p.Points)-1)
	}
	return p.Points[gd.Clampi(index, 0, len(p.Points)-1)]
}

func (p *Path) GetSegment(aIndex, bIndex int) [2]vector2.Vector2 {
	return [2]vector2.Vector2{p.GetPoint(aIndex), p.GetPoint(bIndex)}
}

func (p *Path) GetNearestPoint(fromPosition vector2.Vector2) (vector2.Vector2, int, int) {
	a, b := p.GetNearestSegment(fromPosition)
	return geometry2d.GetClosestPointToSegment(fromPosition, p.GetSegment(a, b)), a, b
}

func (p *Path) GetNearestSegment(fromPosition vector2.Vector2) (int, int) {
	var nSeg [2]int
	nDist := math.Inf(1)
	for i := 1; i < len(p.Points); i++ {
		seg := p.GetSegment(i, i-1)
		dist := geometry2d.GetDistanceSquaredToSegment(fromPosition, seg)
		if dist < nDist {
			nSeg = [2]int{i - 1, i}
			nDist = dist
		}
	}
	return nSeg[0], nSeg[1]
}

func (p *Path) Offset(delta float64, allowInversion bool) error {
	pts := geometry2d.OffsetPolygon(p.Points, delta, geometry2d.JoinTypeMiter)[0]
	if geometry2d.IsPolygonClockwise(pts) && !allowInversion {
		return errors.New("path inverted")
	}
	p.Points = pts
	return nil
}

func (p *Path) GetDistance(fromIndex, toIndex int, direction Direction) float64 {
	_, dist := p.WalkToIndex(fromIndex, toIndex, direction)
	return dist
}

func (p *Path) WalkToEnd(start vector2.Vector2, direction Direction) ([]int, vector2.Vector2) {
	if len(p.Points) < 2 {
		return nil, vector2.Vector2{}
	}

	var indices []int
	_, a, b := p.GetNearestPoint(start)
	end := start
	if p.Closed {
		indices, _ = p.WalkToIndex(b, a, direction)
	} else {
		indices, _ = p.WalkToIndex(b, len(p.Points)-1, direction)
		end = p.Points[len(p.Points)-1]
	}
	return indices, end
}

func (p *Path) WalkToIndex(fromIndex, toIndex int, direction Direction) ([]int, float64) {
	if len(p.Points) < 2 {
		return nil, 0.0
	}

	indices := []int{fromIndex}
	currentIdx := fromIndex
	distance := 0.0

	for currentIdx != toIndex {
		nextIdx := p.Next(currentIdx, direction)
		if !p.Closed && (nextIdx < 0 || nextIdx > len(p.Points)-1) {
			break
		}
		indices = append(indices, currentIdx)
		distance += p.Points[currentIdx].DistanceTo(p.Points[nextIdx])
	}

	return indices, distance
}

func (p *Path) Walk(fromPosition vector2.Vector2, distance float64, direction Direction) ([]int, vector2.Vector2) {
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
			endPosition := vector2.Vector2{
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
		return gd.Wrapi(index, 0, len(p.Points)-1)
	} else {
		return gd.Clampi(index, 0, len(p.Points)-1)
	}
}

func (p *Path) Previous(index int, direction Direction) int {
	if direction == DirectionForward {
		index -= 1
	} else {
		index += 1
	}
	if p.Closed {
		return gd.Wrapi(index, 0, len(p.Points)-1)
	} else {
		return gd.Clampi(index, 0, len(p.Points)-1)
	}
}
