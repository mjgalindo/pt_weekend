package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/mjgalindo/pt_weekend/cam"
	"github.com/mjgalindo/pt_weekend/geo"
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

func color(r ray.Ray, world geo.Hitable) vec.Vec3 {
	var rec geo.HitResult
	if world.Hit(r, 0.0, math.MaxFloat32, &rec) {
		rec.Normal = rec.Normal.MakeUnit()
		return vec.MulSingle(
			vec.Make(rec.Normal.X()+1, rec.Normal.Y()+1, rec.Normal.Z()+1), 0.5)
	}
	unitDir := r.Direction.MakeUnit()
	t := 0.5 * (unitDir.Y() + 1.0)
	return vec.Sum(vec.MulSingle(vec.Make(1.0, 1.0, 1.0), 1.0-t),
		vec.MulSingle(vec.Make(0.5, 0.7, 1.0), t))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need a file parameter")
		os.Exit(1)
	}
	f, err := os.Create(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	width := 400
	height := 200

	// Setup the scene
	world := geo.MakeList(
		geo.Sphere{Position: vec.Make(0, 0, -1), Radius: 0.5},
		geo.Sphere{Position: vec.Make(0, -100.5, -1), Radius: 100})

	camera := cam.Default()
	nSamples := 8

	fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			col := vec.Make(0, 0, 0)
			for s := 0; s < nSamples; s++ {
				u := (float32(x) + rand.Float32()) / float32(width)
				v := (float32(y) + rand.Float32()) / float32(height)
				ray := camera.GetRay(u, v)
				col = vec.Sum(col, color(ray, world))
			}
			col = vec.DivSingle(col, float32(nSamples))
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
		}
	}
}
