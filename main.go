package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"raytracer/sphere"
)

func main() {
	spheres := []*sphere.Sphere{}
	spheres = append(spheres, &sphere.Sphere{
		Center: []float64{5, 1, -12},
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	spheres = append(spheres, &sphere.Sphere{
		Center: []float64{-1.5, -2, -18},
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	spheres = append(spheres, &sphere.Sphere{
		Center: []float64{1.5, -0.5, -22},
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	spheres = append(spheres, &sphere.Sphere{
		Center: []float64{-5, 2, -16},
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	render(spheres)
}

func render(spheres []*sphere.Sphere) {
	width, height := 1024, 768
	framebuf := make([]uint8, width*height*3) // 3 is RBG
	rayOrig := []float64{0, 0, 0}             // camera position
	fov := math.Pi / 2                        // field of view
	startP, unit := GetCanvasStartPointAndUnit(float64(width), float64(height), fov)
	for i := 0; i < height; i++ {
		bgR := uint8(float64(i) / float64(height) * 255)
		for j := 0; j < width; j++ {
			// ray = each vertex on the canvas - camera poision
			rayDir := []float64{
				startP[0] + float64(j)*unit - rayOrig[0],
				startP[1] - float64(i)*unit - rayOrig[1],
				startP[2] - rayOrig[2],
			}
			bgG := uint8(float64(j) / float64(width) * 255)
			color := []uint8{bgR, bgG, 0} // background color
			mindis := math.MaxFloat64
			for _, sphere := range spheres {
				if inter, dis := sphere.IntersectRay(rayOrig, rayDir); inter {
					if dis < mindis {
						color = sphere.Color // closest sphere color
						mindis = dis
					}
				}
			}
			copy(framebuf[(i*width+j)*3:], color)
		}
	}
	file, err := os.Create("out.ppm")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(fmt.Sprintf("P6\n%d %d\n255\n", width, height))
	file.Write(framebuf)
	file.Close()
}

// Start point is the up-left-most position of the canvas
// Unit is the length to move when iterating through pixels
func GetCanvasStartPointAndUnit(width, height, fov float64) ([]float64, float64) {
	if height > width {
		startP := []float64{
			-(math.Tan(fov/2) * (float64(width) - 1) / float64(width)) * float64(width) / float64(height),
			math.Tan(fov/2) * (float64(height) - 1) / float64(height),
			-1,
		}
		unit := math.Tan(fov/2) / float64(height) * 2
		return startP, unit
	}

	startP := []float64{
		-(math.Tan(fov/2) * (float64(width) - 1) / float64(width)),
		math.Tan(fov/2) * (float64(height) - 1) / float64(height) * float64(height) / float64(width),
		-1,
	}
	unit := math.Tan(fov/2) / float64(width) * 2
	return startP, unit
}
