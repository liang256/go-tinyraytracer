package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"raytracer/light"
	"raytracer/material"
	"raytracer/sphere"

	"github.com/deeean/go-vector/vector3"
)

func main() {
	blue := &material.Material{
		Color:     vector3.New(0, 255, 255),
		SpecColor: vector3.New(255, 255, 255),
		Albedo:    vector3.New(0.8, 0.1, 0.0),
		SpecExpo:  1.2,
	}
	red := &material.Material{
		Color:     vector3.New(230, 40, 40),
		SpecColor: vector3.New(255, 255, 255),
		Albedo:    vector3.New(0.4, 0.4, 0.2),
		SpecExpo:  100,
	}
	black := &material.Material{
		Color:     vector3.New(0, 0, 0),
		SpecColor: vector3.New(255, 255, 255),
		Albedo:    vector3.New(0.2, 0.8, 0.9),
		SpecExpo:  1024,
	}
	white := &material.Material{
		Color:     vector3.New(250, 250, 250),
		SpecColor: vector3.New(255, 255, 255),
		Albedo:    vector3.New(0.6, 0.2, 0.0),
		SpecExpo:  800,
	}
	spheres := []*sphere.Sphere{}
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(5, 1, -12),
		Radius:   3,
		Material: red,
	})
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(-1.5, -2, -18),
		Radius:   3,
		Material: black,
	})
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(1.5, -0.5, -22),
		Radius:   3,
		Material: white,
	})
	spheres = append(spheres, &sphere.Sphere{
		Center:   vector3.New(-5, 2, -16),
		Radius:   3,
		Material: blue,
	})

	lights := []*light.Light{}
	lights = append(lights, &light.Light{
		Center:    vector3.New(-20, 0, -10),
		Intensity: 0.4,
	})
	lights = append(lights, &light.Light{
		Center:    vector3.New(40, 15, 20),
		Intensity: 0.4,
	})
	lights = append(lights, &light.Light{
		Center:    vector3.New(2, 30, 20),
		Intensity: 0.1,
	})
	lights = append(lights, &light.Light{
		Center:    vector3.New(-7, 30, 10),
		Intensity: 0.5,
	})
	lights = append(lights, &light.Light{
		Center:    vector3.New(0, 30, 30),
		Intensity: 0.1,
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
		for j := 0; j < width; j++ {
			// ray = each vertex on the canvas - camera poision
			rayDir := startP.AddScalars(float64(j)*unit, -(float64(i) * unit), 0).Sub(rayOrig)
			color := vec3ToUint8(castRay(rayOrig, rayDir, spheres, lights, 0))
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
func castRay(rayOrig, rayDir *vector3.Vector3, spheres []*sphere.Sphere, lights []*light.Light, depth int) *vector3.Vector3 {
	if depth > 4 {
		return vector3.New(0, 0, 0)
	}
	rayDir = rayDir.Normalize()
	mindis := math.MaxFloat64
	hitSphereId := -1
	for i, sphere := range spheres {
		if inter, dis := sphere.IntersectRay(rayOrig, rayDir); inter && dis < mindis {
			mindis, hitSphereId = dis, i
		}
	}
	if mindis == math.MaxFloat64 {
		return vector3.New(200, 200, 0) // bg color
	}
	hitPoint := rayOrig.Add(rayDir.MulScalar(mindis))
	hitN := hitPoint.Sub(spheres[hitSphereId].Center).Normalize()

	reflectColor := castRay(hitPoint, reflect(rayDir, hitN).Normalize(), spheres, lights, depth+1)
	diffuseItensity, specIntensity := 0.0, 0.0
	for _, l := range lights {
		lightDir := hitPoint.Sub(l.Center).Normalize()
		v := hitN.Dot(lightDir)
		if v < 0 {
			// from hit point, check if there is an object between the point and the light
			min, hitId := math.MaxFloat64, -1
			for j := range spheres {
				hit, dis := spheres[j].IntersectRay(l.Center, lightDir)
				if hit && dis < min {
					min, hitId = dis, j
				}
			}
			if hitId == hitSphereId {
				m := spheres[hitSphereId].Material
				specIntensity += math.Pow(math.Max(0, reflect(lightDir, hitN).Dot(lightDir.MulScalar(-1))), m.SpecExpo) * l.Intensity
				diffuseItensity += -v * l.Intensity
			}
		}
	}
	m := spheres[hitSphereId].Material
	return m.Color.MulScalar(diffuseItensity * m.Albedo.X).Add(m.SpecColor.MulScalar(specIntensity * m.Albedo.Y)).Add(reflectColor.MulScalar(m.Albedo.Z))
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
