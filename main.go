package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"raytracer/light"
	"raytracer/sphere"

	"github.com/deeean/go-vector/vector3"
)

func main() {
	defer duration(track("main"))
	spheres := []*sphere.Sphere{}
	spheres = append(spheres, &sphere.Sphere{
		Center: vector3.New(5, 1, -12),
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	spheres = append(spheres, &sphere.Sphere{
		Center: vector3.New(-1.5, -2, -18),
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	spheres = append(spheres, &sphere.Sphere{
		Center: vector3.New(1.5, -0.5, -22),
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	spheres = append(spheres, &sphere.Sphere{
		Center: vector3.New(-5, 2, -16),
		Radius: 3,
		Color:  []uint8{0, 255, 255},
	})
	lights := []*light.Light{}
	lights = append(lights, &light.Light{
		Center:    vector3.New(-20, 0, -10),
		Intensity: 0.1,
	})
	lights = append(lights, &light.Light{
		Center:    vector3.New(40, 15, 20),
		Intensity: 0.9,
	})
	render(spheres, lights)
}

func render(spheres []*sphere.Sphere, lights []*light.Light) {
	width, height := 1024, 768
	framebuf := make([]uint8, width*height*3) // 3 is RBG
	rayOrig := vector3.New(0, 0, 0)           // camera position
	fov := math.Pi / 2                        // field of view
	startP, unit := GetCanvasStartPointAndUnit(float64(width), float64(height), fov)
	for i := 0; i < height; i++ {
		bgR := uint8(float64(i) / float64(height) * 255)
		for j := 0; j < width; j++ {
			// ray = each vertex on the canvas - camera poision
			rayDir := startP.AddScalars(float64(j)*unit, -(float64(i) * unit), 0).Sub(rayOrig)
			bgG := uint8(float64(j) / float64(width) * 255)
			color := []uint8{bgR, bgG, 0} // default color is background
			castRay(rayOrig, rayDir, spheres, lights, &color)
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

// If the ray hit a geo, cast its color to the color-array-ptr
func castRay(rayOrig, rayDir *vector3.Vector3, spheres []*sphere.Sphere, lights []*light.Light, color *[]uint8) {
			mindis := math.MaxFloat64
			for _, sphere := range spheres {
		if inter, dis := sphere.IntersectRay(rayOrig, rayDir); inter && dis < mindis {
			// the point that the ray hits on the sphere
				pOnSphere := rayOrig.Add(rayDir.Normalize().MulScalar(dis))
			// normal of the point
				normal := pOnSphere.Sub(sphere.Center).Normalize()
				diffuseItensity := 0.0
			// the light angle is more close to the normal
			// or the more lights causes a greater intensity
				for _, l := range lights {
					lightDir := l.Center.Sub(pOnSphere).Normalize()
					diffuseItensity += math.Max(lightDir.Dot(normal)*l.Intensity, 0)
						}
				diffuseItensity = math.Min(diffuseItensity, 1)
				*color = []uint8{
					uint8(float64(sphere.Color[0]) * diffuseItensity),
					uint8(float64(sphere.Color[1]) * diffuseItensity),
					uint8(float64(sphere.Color[2]) * diffuseItensity),
						}
						mindis = dis
					}
				}
			}
}

// Start point is the up-left-most position of the canvas
// Unit is the length to move when iterating through pixels
func GetCanvasStartPointAndUnit(width, height, fov float64) (*vector3.Vector3, float64) {
	startP := vector3.New(
		-(math.Tan(fov/2) * (width - 1) / width),
		math.Tan(fov/2)*(height-1)/height,
		-1,
	)
	if height > width {
		startP = startP.MulScalars(width/height, 1, 1)
	} else if height < width {
		startP = startP.MulScalars(1, height/width, 1)
		}
	// unit = mmath.Tan(fov/2) / width * 2
	return startP, math.Tan(fov/2) / width * 2
	}

	startP := []float64{
		-(math.Tan(fov/2) * (float64(width) - 1) / float64(width)),
		math.Tan(fov/2) * (float64(height) - 1) / float64(height) * float64(height) / float64(width),
		-1,
	}
	unit := math.Tan(fov/2) / float64(width) * 2
	return startP, unit
}
