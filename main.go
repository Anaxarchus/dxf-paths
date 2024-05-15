package main

import (
	paths "dxf-paths/pkg/path"
	"fmt"

	"github.com/Anaxarchus/gdscript-libs/pkg/mathgd"
)

func main() {
	pts := []paths.VectorB{
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 0, Y: 5.5},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 0, Y: 0},
		},
		{
			B:       0.19891,
			Vector2: mathgd.Vector2{X: 17.58581, Y: 0.0},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 18.29291, Y: 0.29289},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 19.125, Y: 1.125},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 19.125, Y: 12.425},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 17.45, Y: 12.425},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 17.45, Y: 12.625},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 15.55, Y: 12.625},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 15.55, Y: 12.425},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 13.625, Y: 12.425},
		},
		{
			B:       0,
			Vector2: mathgd.Vector2{X: 13.625, Y: 5.5},
		},
	}
	path := paths.NewPathFromBulge(pts, true, 0.1)

	snippets := path.GetArcSnippets()
	fmt.Println("snippets: ", snippets)
	for i := range snippets {
		fmt.Println("arc: ", snippets[i].Arc)
	}
}

func printToGodot(points []mathgd.Vector2) {
	fmt.Print("\nvar points:PackedVector2Array = [")
	for i := range points {
		fmt.Printf("\n	Vector2(%f,%f),", points[i].X, points[i].Y)
	}
	fmt.Print("\n]")
}
