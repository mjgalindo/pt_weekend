package main

import (
	"fmt"

	"github.com/mjgalindo/pt_weekend/vec"
)

func main() {
	width := 200
	height := 100
	fmt.Printf("P3\n%d %d\n255\n", width, height)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			col := vec.Make(float32(x)/float32(width), float32(y)/float32(height), 0.2)
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			fmt.Printf("%d %d %d\n", ir, ig, ib)
		}
	}
}
