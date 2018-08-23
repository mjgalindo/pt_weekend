package main

import (
	"fmt"

	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

func color(r ray.Ray) vec.Vec3 {
	if hitSphere(vec.Make(0.0, 0.0, -1.0), 0.5, r) {
		return vec.Make(1.0, 0.0, 0.0)
	}
	unitDir := r.Direction.MakeUnit()
	t := 0.5 * (unitDir.Y() + 1.0)
	return vec.Sum(vec.MulSingle(vec.Make(1.0, 1.0, 1.0), 1.0-t), vec.MulSingle(vec.Make(0.5, 0.8, 1.0), t))
}

func hitSphere(center vec.Vec3, radius float32, r ray.Ray) bool {
	oc := vec.Sub(r.Origin, center)
	a := vec.Dot(r.Direction, r.Direction)
	b := 2.0 * vec.Dot(oc, r.Direction)
	c := vec.Dot(oc, oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant > 0
}

func main() {
	width := 200
	height := 100
	fmt.Printf("P3\n%d %d\n255\n", width, height)
	lowerLeftCorner := vec.Make(-2.0, -1.0, -1.0)
	horizontal := vec.Make(4.0, 0.0, 0.0)
	vertical := vec.Make(0.0, 2.0, 0.0)
	origin := vec.Make(0.0, 0.0, 0.0)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			u := float32(x) / float32(width)
			v := float32(y) / float32(height)
			ray := ray.Make(origin, vec.Sum(lowerLeftCorner, vec.MulSingle(horizontal, u), vec.MulSingle(vertical, v)))
			col := color(ray)
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}
