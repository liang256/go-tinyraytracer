package sphere

import (
	"raytracer/material"
	"testing"

	"github.com/deeean/go-vector/vector3"
)

func TestIntersectRay(t *testing.T) {
	rayOrig, rayDir := vector3.New(0, 0, 0), vector3.New(0, 0, -1)
	m := &material.Material{}
	// sphere 1
	sphere := Sphere{vector3.New(4, 3, -16), 5, m}
	if intersect, _ := sphere.IntersectRay(rayOrig, rayDir); !intersect {
		t.Error()
	}
	// sphere 2
	sphere = Sphere{vector3.New(4, 10, -16), 3, m}
	if intersect, _ := sphere.IntersectRay(rayOrig, rayDir); intersect {
		t.Error()
	}
	// sphere 3
	sphere = Sphere{vector3.New(0, 0, -16), 3, m}
	if _, dis := sphere.IntersectRay(rayOrig, rayDir); dis != 13 {
		t.Error()
	}
}
