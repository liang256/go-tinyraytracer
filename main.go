package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"raytracer/light"
	"raytracer/material"
	"raytracer/plane"
	"raytracer/sphere"

	"github.com/deeean/go-vector/vector3"
)

type Renderable interface {
	IntersectRay(orig, dir *vector3.Vector3) (bool, float64)
	GetCenter() *vector3.Vector3
	GetMaterial() *material.Material
}

func main() {
	blue := &material.Material{
		Color:      vector3.New(0, 255, 255),
		SpecColor:  vector3.New(255, 255, 255),
		Albedo:     vector3.New(0.8, 0.1, 0.0),
		SpecExpo:   5,
		Refractive: 1,
	}
	red := &material.Material{
		Color:      vector3.New(230, 40, 40),
		SpecColor:  vector3.New(255, 255, 255),
		Albedo:     vector3.New(0.5, 0.4, 0.1),
		SpecExpo:   100,
		Refractive: 1,
	}
	black := &material.Material{
		Color:      vector3.New(0, 0, 0),
		SpecColor:  vector3.New(255, 255, 255),
		Albedo:     vector3.New(0.2, 0.8, 0.9),
		SpecExpo:   1024,
		Refractive: 1.2,
	}
	white := &material.Material{
		Color:      vector3.New(250, 250, 250),
		SpecColor:  vector3.New(255, 255, 255),
		Albedo:     vector3.New(0.6, 0.2, 0.1),
		SpecExpo:   800,
		Refractive: 1,
	}
	geos := []Renderable{}
	geos = append(geos, &sphere.Sphere{
		Center:   vector3.New(5, 1, -12),
		Radius:   3,
		Material: red,
	})
	geos = append(geos, &sphere.Sphere{
		Center:   vector3.New(-1.5, -2, -18),
		Radius:   3,
		Material: black,
	})
	geos = append(geos, &sphere.Sphere{
		Center:   vector3.New(1.5, -0.5, -22),
		Radius:   3,
		Material: white,
	})
	geos = append(geos, &sphere.Sphere{
		Center:   vector3.New(-5, 2, -16),
		Radius:   3,
		Material: blue,
	})
	geos = append(geos, &plane.Disc{
		Center:   vector3.New(0, -6, -16),
		Normal:   vector3.New(0.1, 1, -0.2).Normalize(),
		Radius:   5,
		Material: blue,
	})

	lights := []*light.Light{}
	lights = append(lights, &light.Light{
		Center:    vector3.New(0, 100, -10),
		Intensity: 1,
	})
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
	lights = append(lights, &light.Light{
		Center:    vector3.New(0, -30, 30),
		Intensity: 0.1,
	})
	render(geos, lights)
}

func render(geos []Renderable, lights []*light.Light) {
	width, height := 1024, 768
	framebuf := make([]uint8, width*height*3) // 3 is RBG
	rayOrig := vector3.New(0, 0, 0)           // camera position
	fov := math.Pi / 2                        // field of view
	startP, unit := GetCanvasStartPointAndUnit(float64(width), float64(height), fov)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			// ray = each vertex on the canvas - camera poision
			rayDir := startP.AddScalars(float64(j)*unit, -(float64(i) * unit), 0).Sub(rayOrig)
			color := castRay(rayOrig, rayDir, geos, lights, 0)
			color = color.MulScalar(0.95).Add(vector3.New(200, 200, 0).MulScalar(0.05))
			copy(framebuf[(i*width+j)*3:], vec3ToUint8(color))
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

func refract(in, normal *vector3.Vector3, refractive float64) *vector3.Vector3 {
	cosi := math.Min(1, math.Max(0, in.Dot(normal)))
	etai := 1.0
	etat := refractive
	if cosi < 0 {
		cosi *= -cosi
		etat, etai = etai, etat
		normal = normal.MulScalar(-1)
	}
	eta := etai / etat
	k := 1 - eta*eta*(1-cosi*cosi)
	if k < 0 {
		return vector3.New(0, 0, 0)
	}
	return in.MulScalar(eta).Add(normal.MulScalar(cosi).SubScalar(math.Sqrt(k)))
}

// If the ray hit a geo, cast its color to the color-array-ptr
func castRay(rayOrig, rayDir *vector3.Vector3, geos []Renderable, lights []*light.Light, depth int) *vector3.Vector3 {
	if depth > 4 {
		return vector3.New(0, 0, 0)
	}
	rayDir = rayDir.Normalize()
	mindis := math.MaxFloat64
	hitGeoId := -1
	for i, geo := range geos {
		if inter, dis := geo.IntersectRay(rayOrig, rayDir); inter && dis < mindis {
			mindis, hitGeoId = dis, i
		}
	}
	if mindis == math.MaxFloat64 {
		return vector3.New(200, 200, 0) // bg color
	}
	// return geos[hitGeoId].GetMaterial().Color
	hitPoint := rayOrig.Add(rayDir.Normalize().MulScalar(mindis))
	var hitN *vector3.Vector3
	switch geos[hitGeoId].(type) {
	case *sphere.Sphere:
		hitN = hitPoint.Sub(geos[hitGeoId].GetCenter()).Normalize()
	case *plane.Disc:
		hitN = geos[hitGeoId].(*plane.Disc).Normal.Normalize()
	}
	m := geos[hitGeoId].GetMaterial()

	refractDir := refract(rayDir, hitN, m.Refractive).Normalize()
	refractOrig := hitPoint.Add(hitN.MulScalar(1e-3))
	if refractDir.Dot(hitN) < 0 {
		refractOrig = hitPoint.Sub(hitN.MulScalar(1e-3))
	}
	reflectColor := castRay(hitPoint, reflect(rayDir, hitN).Normalize(), geos, lights, depth+1)
	refractColor := castRay(refractOrig, refractDir, geos, lights, depth+1)
	diffuseItensity, specIntensity := 0.0, 0.0
	for _, l := range lights {
		lightDir := hitPoint.Sub(l.Center).Normalize()
		v := hitN.Dot(lightDir)
		if v < 0 {
			// from hit point, check if there is an object between the point and the light
			min, hitId := math.MaxFloat64, -1
			for j := range geos {
				hit, dis := geos[j].IntersectRay(l.Center, lightDir)
				if hit && dis < min {
					min, hitId = dis, j
				}
			}
			if hitId == hitGeoId {
				specIntensity += math.Pow(math.Max(0, reflect(lightDir, hitN).Dot(lightDir.MulScalar(-1))), m.SpecExpo) * l.Intensity
				diffuseItensity += -v * l.Intensity
			}
		}
	}
	return m.Color.MulScalar(diffuseItensity * m.Albedo.X).Add(m.SpecColor.MulScalar(specIntensity * m.Albedo.Y)).Add(reflectColor.MulScalar(m.Albedo.Z)).Add(refractColor.MulScalar(m.Albedo.Z))
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
