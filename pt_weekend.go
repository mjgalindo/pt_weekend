package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"gopkg.in/cheggaaa/pb.v1"

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
			}
			return vec.Make(0, 0, 0)
		}
	}
	unitDir := r.Direction.MakeUnit()
	t := 0.5 * (unitDir.Y() + 1.0)
	return vec.Sum(vec.MulSingle(vec.Make(1.0, 1.0, 1.0), 1.0-t),
		vec.MulSingle(vec.Make(0.5, 0.7, 1.0), t))
}

func Save(image *[][]vec.Vec3, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "P3\n%d %d\n255\n", len((*image)[0]), len(*image))
	for y := len(*image) - 1; y >= 0; y-- {
		for x := range (*image)[y] {
			ir := int(255.99 * (*image)[y][x].R())
			ig := int(255.99 * (*image)[y][x].G())
			ib := int(255.99 * (*image)[y][x].B())
			fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
		}
	}
}
func renderScene(outfile string) {
	width := 800 / 2
	height := 400 / 2

	// Setup the scene
	world := randomScene()
	lookFrom := vec.Make(3, 2, 1)
	lookAt := vec.Make(0, 0, -1)
	distToFocus := vec.Sub(lookFrom, lookAt).Length()
	aperture := float32(0.25)
	camera := cam.Make(lookFrom, lookAt, vec.Make(0, 1, 0),
		90, float32(width)/float32(height), aperture, distToFocus)
	nSamples := 64
	image := make([][]vec.Vec3, height)
	for i := range image {
		image[i] = make([]vec.Vec3, width)
	}

	renderWorker := func(workQueue chan int, finish chan bool) {
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
						image[row][x] = vec.Sum(image[row][x], color(ray, world, 0))
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
	nWorkers := 8
	for i := 0; i < nWorkers; i++ {
		go renderWorker(workQueue, finish)
	}
	bar := pb.StartNew(height)
	defer bar.FinishPrint("Done!")
	for y := height - 1; y >= 0; y-- {
		select {
		case workQueue <- y:
			bar.Increment()
		}
	}
	for i := 0; i < nWorkers; i++ {
		finish <- true
	}

	Save(&image, outfile)
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need a file parameter")
		os.Exit(1)
	}
	renderScene(os.Args[1])
}

func randomScene() *geo.HitableList {
	list := geo.MakeList()
	list.Add(geo.MakeSphere(vec.Make(0, -1000, 0), 1000, mat.Lambertian(vec.Make(0.5, 0.5, 0.5))))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := rand.Float32()
			center := vec.Make(float32(a)+0.9*rand.Float32(), 0.2, float32(b)+0.9*rand.Float32())
			if vec.Sub(center, vec.Make(4, 0.2, 0)).Length() > 0.9 {
				if chooseMat < 0.8 { // Diffuse
					albedo := vec.Make(rand.Float32()*rand.Float32(), rand.Float32()*rand.Float32(), rand.Float32()*rand.Float32())
					list.Add(geo.MakeSphere(center, 0.2, mat.Lambertian(albedo)))
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
	return &list
}
