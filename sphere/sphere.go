package sphere

import (
	"math"
	"raytracer/material"

	"github.com/deeean/go-vector/vector3"
)

type Sphere struct {
	Center   *vector3.Vector3
	Radius   float64
	Material *material.Material
}

// Returns isIntersect, distant from point to camera, and normal of the point
func (s Sphere) IntersectRay(orig, dir *vector3.Vector3) (bool, float64) {
	// if the sphere not at the side of ray direction return false
	dir = dir.Normalize()
	t := s.Center.Sub(orig).Dot(dir)
	p := orig.Add(dir.MulScalar(t))
	y := s.Center.Distance(p)
	if y > s.Radius {
		return false, 0.0
	} else if orig.Distance(s.Center) < s.Radius {
		return false, 0.0 // case: ray-orig in the sphere
	} else if t < 0 {
		return false, 0.0 // case: the sphere is not at the ray-dir side
	}
	// caculate the distance between the closest point on this sphere and the camera
	x := math.Sqrt(s.Radius*s.Radius - y*y)
	return true, t - x
}
