package paths

import "github.com/Anaxarchus/gdscript-libs/pkg/mathgd"

type Snippet struct {
	Indices []int
	Path    *Path
	Start   mathgd.Vector2
	End     mathgd.Vector2
}

type ArcSnippet struct {
	Snippet
	Arc *Arc
}

func NewArcSnippet(start, end int, path *Path) *ArcSnippet {
	indices := []int{}
	points := []mathgd.Vector2{}
	n := len(path.Points)

	for i := start; ; i = (i + 1) % n {
		indices = append(indices, i)
		points = append(points, path.Points[i])
		if i == end {
			break
		}
	}
	return &ArcSnippet{
		Snippet: Snippet{
			Path:    path,
			Indices: indices,
			Start:   points[0],
			End:     points[len(points)-1],
		},
		Arc: NewArcFromArcPoints(points),
	}
}

func (a *ArcSnippet) GetArcSection(iFrom, iTo int) *Arc {
	return &Arc{
		Origin:     a.Arc.Origin,
		Radius:     a.Arc.Radius,
		StartAngle: a.Arc.AngleToPoint(a.Path.GetPoint(iFrom)),
		EndAngle:   a.Arc.AngleToPoint(a.Path.GetPoint(iTo)),
	}
}
