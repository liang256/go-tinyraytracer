package sphere

import (
	"math"
)

type Sphere struct {
	Center []float64
	Radius float64
	Color  []uint8 // RGB
}

func (s Sphere) IntersectRay(orig, dir []float64) (bool, float64) {
	// if the sphere not at the side of ray direction return false
	origCenterVec := SubtractMatrix(s.Center, orig)
	ray := AddMatrix(orig, dir)
	dotProd := DotProduct(origCenterVec, ray)
	if dotProd < 0 {
		return false, 0.0
	}
	// if the ray orig in the sphere, return false
	if VecLen(origCenterVec) <= s.Radius {
		return false, 0.0
	}
	// get the projected center on the ray, pc
	// length of projection
	projLen := dotProd / VecLen(ray)
	// unit-ray-vector
	unitRay := ScaleMatrix(ray, 1.0/VecLen(ray))
	// pc = orig + unitray * projectlength
	pc := AddMatrix(orig, ScaleMatrix(unitRay, projLen))
	// if |pc - c| <= radius, return true
	if res := VecLen(SubtractMatrix(pc, s.Center)); res > s.Radius {
		return false, 0.0
	}
	// caculate the distance between the closest point on this sphere and the camera
	dis := projLen - math.Sqrt(math.Pow(s.Radius, 2)-math.Pow(VecLen(origCenterVec), 2)+math.Pow(projLen, 2))
	return true, dis
}

func SubtractMatrix(a, b []float64) []float64 {
	res := make([]float64, len(a))
	for i := range a {
		res[i] = a[i] - b[i]
	}
	return res
}

func AddMatrix(a, b []float64) []float64 {
	res := make([]float64, len(a))
	for i := range a {
		res[i] = a[i] + b[i]
	}
	return res
}

func ScaleMatrix(vec []float64, s float64) []float64 {
	res := make([]float64, len(vec))
	for i := range vec {
		res[i] = vec[i] * s
	}
	return res
}

func DotProduct(a, b []float64) float64 {
	res := 0.0
	for i := range a {
		res += a[i] * b[i]
	}
	return res
}

func VecLen(vec []float64) float64 {
	sum := 0.0
	for _, val := range vec {
		sum += val * val
	}
	return math.Sqrt(sum)
}
