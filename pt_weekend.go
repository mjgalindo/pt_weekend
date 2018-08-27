package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/mjgalindo/pt_weekend/cam"
	"github.com/mjgalindo/pt_weekend/geo"
	"github.com/mjgalindo/pt_weekend/mat"
	"github.com/mjgalindo/pt_weekend/ray"
	"github.com/mjgalindo/pt_weekend/vec"
)

func colorAt(r ray.Ray, world geo.Hitable, depth int) vec.Vec3 {
	var rec geo.HitRecord
	if world.Hit(r, 0.001, math.MaxFloat32, &rec) {
		if depth < 50 {
			absorbed, attenuation, scattered := rec.Scatter(r, rec)
			if !absorbed {
				return vec.Mul(attenuation, colorAt(scattered, world, depth+1))
			}
			return vec.Make(0, 0, 0)
		}
	}
	unitDir := r.Direction.MakeUnit()
	t := 0.5 * (unitDir.Y() + 1.0)
	return vec.Sum(vec.MulSingle(vec.Make(1.0, 1.0, 1.0), 1.0-t),
		vec.MulSingle(vec.Make(0.5, 0.7, 1.0), t))
}

func save(imbuff *[][]vec.Vec3, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	width, height := len((*imbuff)[0]), len(*imbuff)
	pngImage := image.NewRGBA(image.Rectangle{Max: image.Point{X: width, Y: height}, Min: image.Point{X: 0, Y: 0}})
	for y := len(*imbuff) - 1; y >= 0; y-- {
		for x := range (*imbuff)[y] {
			ir := uint8(255.99 * (*imbuff)[y][x].R())
			ig := uint8(255.99 * (*imbuff)[y][x].G())
			ib := uint8(255.99 * (*imbuff)[y][x].B())
			pngImage.SetRGBA(x, height-1-y, color.RGBA{R: ir, G: ig, B: ib, A: 255})
		}
	}
	encoder := png.Encoder{}
	encoder.Encode(f, pngImage)
}

func renderScene(outfile string) {
	rand.Seed(int64(time.Now().Nanosecond()))
	width := 1600
	height := 800

	// Setup the scene
	world := randomScene()
	lookFrom := vec.Make(13, 3.5, 3)
	lookAt := vec.Make(0, 0, 0)
	distToFocus := float32(1.0)
	aperture := float32(0.0)
	camera := cam.Make(lookFrom, lookAt, vec.Make(0, 1, 0),
		20, float32(width)/float32(height), aperture, distToFocus, 0.0, 1.0)
	nSamples := 128

	// Initialize a 2D matrix of Vec3 to store the radiance at each pixel
	// Could be int8 directly but we may want to do something with the extra
	// precision (?)
	image := make([][]vec.Vec3, height)
	for i := range image {
		image[i] = make([]vec.Vec3, width)
	}

	worker := func(workQueue chan int, finish chan bool) {
		mustEnd := false
		for !mustEnd {
			select {
			case row := <-workQueue:
				for x := 0; x < width; x++ {
					image[row][x] = vec.Make(0, 0, 0)
					for s := 0; s < nSamples; s++ {
						u := (float32(x) + rand.Float32()) / float32(width)
						v := (float32(row) + rand.Float32()) / float32(height)
						ray := camera.GetRay(u, v)
						image[row][x] = vec.Sum(image[row][x], colorAt(ray, world, 0))
					}
					image[row][x] = vec.DivSingle(image[row][x], float32(nSamples))
					image[row][x] = vec.Make(float32(math.Sqrt(float64(image[row][x].X()))),
						float32(math.Sqrt(float64(image[row][x].Y()))),
						float32(math.Sqrt(float64(image[row][x].Z()))))
				}
			case <-finish:
				mustEnd = true
			}
		}
	}

	workQueue := make(chan int)
	finish := make(chan bool)
	nWorkers := 4
	for i := 0; i < nWorkers; i++ {
		go worker(workQueue, finish)
	}
	bar := pb.StartNew(height)
	// Create a random order for rendering the rows to share the
	// most complicated rows 'equally'
	rowList := make([]int, height)
	for i := range rowList {
		rowList[i] = i
	}
	rand.Shuffle(height, func(i, j int) {
		rowList[i], rowList[j] = rowList[j], rowList[i]
	})
	// Give jobs to the workers (row indices to render)
	for y := height - 1; y >= 0; y-- {
		select {
		case workQueue <- rowList[y]:
			bar.Increment()
		}
	}
	// Stop all workers
	for i := 0; i < nWorkers; i++ {
		finish <- true
	}

	save(&image, outfile)
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need a file parameter")
		os.Exit(1)
	}
	renderScene(os.Args[1])
}

func randomScene() geo.Hitable {
	list := geo.MakeList()
	list.Add(geo.MakeSphere(vec.Make(0, -10000, 0), 10000, mat.Lambertian(vec.Make(0.5, 0.5, 0.5))))
	for a := -100; a < 100; a++ {
		for b := -100; b < 100; b++ {
			chooseMat := rand.Float32()
			center := vec.Make(float32(a)+0.9*rand.Float32(), 0.2, float32(b)+0.9*rand.Float32())
			if vec.Sub(center, vec.Make(4, 0.2, 0)).Length() > 0.9 {
				if chooseMat < 0.8 { // Diffuse
					albedo := vec.Make(rand.Float32()*rand.Float32(), rand.Float32()*rand.Float32(), rand.Float32()*rand.Float32())
					list.Add(geo.MakeMovingSphere(center, vec.Sum(center, vec.Make(0, 0.2*rand.Float32(), 0)), 0.2, 0.0, 1.0, mat.Lambertian(albedo)))
				} else if chooseMat < 0.95 { // Metal
					metalFun := mat.Metal(vec.Make(0.5*(1+rand.Float32()), 0.5*(1+rand.Float32()), 0.5*(1+rand.Float32())), 0.5*rand.Float32())
					list.Add(geo.MakeSphere(center, 0.2, metalFun))
				} else { // Glass
					list.Add(geo.MakeSphere(center, 0.2, mat.Dielectric(1.5)))
				}
			}
		}
	}
	list.Add(geo.MakeSphere(vec.Make(0, 1, 0), 1.0, mat.Dielectric(1.5)))
	list.Add(geo.MakeSphere(vec.Make(-2, 1, 0), 1.0, mat.Lambertian(vec.Make(0.4, 0.2, 0.1))))
	list.Add(geo.MakeSphere(vec.Make(2, 1, 0), 1.0, mat.Metal(vec.Make(0.7, 0.6, 0.5), 0.0)))
	return geo.ConstructBVH(list.Hitables, 0, 1)
}
