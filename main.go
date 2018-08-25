package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/mjgalindo/pt_weekend/cam"
	"github.com/mjgalindo/pt_weekend/geo"
	"github.com/mjgalindo/pt_weekend/mat"
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

func color(r ray.Ray, world geo.Hitable, depth int) vec.Vec3 {
	var rec geo.HitRecord
	if world.Hit(r, 0.001, math.MaxFloat32, &rec) {
		if depth < 50 {
			absorbed, attenuation, scattered := rec.Scatter(r, rec)
			if !absorbed {
				return vec.Mul(attenuation, color(scattered, world, depth+1))
			} else {
				return vec.Make(0, 0, 0)
			}
		}
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
	width := 500 / 2
	height := 250 / 2

	// Setup the scene
	world := geo.MakeList(
		geo.MakeSphere(vec.Make(0, 0, -1), 0.5, mat.Lambertian(vec.Make(0.1, 0.2, 0.5))),
		geo.MakeSphere(vec.Make(0, -100.5, -1), 100, mat.Lambertian(vec.Make(0.8, 0.8, 0.0))),
		geo.MakeSphere(vec.Make(1, 0, -1), 0.5, mat.Mirror(vec.Make(0.8, 0.6, 0.2))),
		geo.MakeSphere(vec.Make(-1, 0, -1), 0.5, mat.Dielectric(1.5)))

	camera := cam.Default()
	nSamples := 64

	fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			col := vec.Make(0, 0, 0)
			for s := 0; s < nSamples; s++ {
				u := (float32(x) + rand.Float32()) / float32(width)
				v := (float32(y) + rand.Float32()) / float32(height)
				ray := camera.GetRay(u, v)
				col = vec.Sum(col, color(ray, world, 0))
			}
			col = vec.DivSingle(col, float32(nSamples))
			col = vec.Make(float32(math.Sqrt(float64(col.X()))),
				float32(math.Sqrt(float64(col.Y()))),
				float32(math.Sqrt(float64(col.Z()))))
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
		}
	}
}
