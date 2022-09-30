package sphere

import "testing"

func TestDotProduct(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{1, 5, 0}
	if res := DotProduct(a, b); res == 11 {
		t.Log("success")
	} else {
		t.Error("fail")
	}
}

func TestVecLen(t *testing.T) {
	vec := []float64{3, 4, 0}
	if res := VecLen(vec); res != 5 {
		t.Error("fail")
	}
}

func TestAddMatrix(t *testing.T) {
	a := []float64{3, 4, 0}
	b := []float64{3, 4, 10}
	res := AddMatrix(a, b)
	for i, val := range []float64{6, 8, 10} {
		if res[i] != val {
			t.Error()
		}
	}
}

func TestSubtractMatrix(t *testing.T) {
	a := []float64{3, 4, 0}
	b := []float64{3, 4, 10}
	res := SubtractMatrix(a, b)
	for i, val := range []float64{0, 0, -10} {
		if res[i] != val {
			t.Error()
		}
	}
}

func TestScaleMatrix(t *testing.T) {
	vec := []float64{3, 4, 0}
	var scale float64 = 0.5
	res := ScaleMatrix(vec, scale)
	for i, val := range []float64{1.5, 2, 0} {
		if res[i] != val {
			t.Error()
		}
	}
}

func TestIntersectRay(t *testing.T) {
	rayOrig, rayDir := []float64{0, 0, 0}, []float64{0, 0, -1}
	// sphere 1
	sphere := Sphere{[]float64{4, 3, -16}, 5, []uint8{0, 0, 0}}
	if intersect, _ := sphere.IntersectRay(rayOrig, rayDir); !intersect {
		t.Error()
	}
	// sphere 2
	sphere = Sphere{[]float64{4, 10, -16}, 3, []uint8{0, 0, 0}}
	if intersect, _ := sphere.IntersectRay(rayOrig, rayDir); intersect {
		t.Error()
	}
	// sphere 3
	sphere = Sphere{[]float64{0, 0, -16}, 3, []uint8{0, 0, 0}}
	if _, dis := sphere.IntersectRay(rayOrig, rayDir); dis != 13 {
		t.Error()
	}
}
