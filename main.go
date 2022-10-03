package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"raytracer/light"
	"raytracer/material"
	"raytracer/sphere"

	"github.com/deeean/go-vector/vector2"
	"github.com/deeean/go-vector/vector3"
)

func main() {
	m := &material.Material{
		Color:     vector3.New(0, 255, 255),
		SpecColor: vector3.New(200, 255, 255),
		Albedo:    vector2.New(0.9, 0.1),
		SpecExpo:  9,
	}
	spheres := []*sphere.Sphere{}
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(5, 1, -12),
		Radius:   3,
		Material: m,
	})
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(-1.5, -2, -18),
		Radius:   3,
		Material: m,
	})
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(1.5, -0.5, -22),
		Radius:   3,
		Material: m,
	})
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(-5, 2, -16),
		Radius:   3,
		Material: m,
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

// Returns the out vector
func reflect(in, normal *vector3.Vector3) *vector3.Vector3 {
	return in.Sub(normal.MulScalar(normal.Dot(in)).MulScalar(2))
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
			diffuseItensity, specIntensity := 0.0, 0.0
			m := sphere.Material
			// the light angle is more close to the normal
			// or the more lights causes a greater intensity
			for _, l := range lights {
				lightDir := l.Center.Sub(pOnSphere).Normalize()
				diffuseItensity += math.Max(lightDir.Dot(normal)*l.Intensity, 0)
				specIntensity += math.Pow(math.Max(0, reflect(lightDir, normal).Dot(rayDir)), m.SpecExpo) * l.Intensity
			}
			diffuseItensity = math.Min(diffuseItensity, 1)
			diffusePass := m.Color.MulScalar(diffuseItensity * m.Albedo.X)
			specPass := m.SpecColor.MulScalar(specIntensity * m.Albedo.Y)
			*color = vec3ToUint8(diffusePass.Add(specPass))
			mindis = dis
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

func vec3ToUint8(v *vector3.Vector3) []uint8 {
	return []uint8{
		uint8(math.Min(math.Max(v.X, 0), 255)),
		uint8(math.Min(math.Max(v.Y, 0), 255)),
		uint8(math.Min(math.Max(v.Z, 0), 255)),
	}
}
