package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"raytracer/sphere"
)

func main() {
	width, height := 1024, 768
	sphere := &sphere.Sphere{
		Center: []float64{-5, 2, -16},
		Radius: 3,
	}
	rayOrig := []float64{0, 0, 0} // camera position
	fov := math.Pi / 2            // field of view
	startP, unit := GetCanvasStartPointAndUnit(float64(width), float64(height), fov)
	file, err := os.Create("out.ppm")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(fmt.Sprintf("P6\n%d %d\n255\n", width, height))
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
			if sphere.IntersectRay(rayOrig, rayDir) {
				color = []uint8{0, 255, 255} // sphere color
			}
			file.Write(color) // RGB
		}
	}
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
