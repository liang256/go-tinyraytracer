package main

import (
	"testing"

	"github.com/deeean/go-vector/vector3"
)

func TestReflect(t *testing.T) {
	inOrig := vector3.New(3, 4, 1)
	in := vector3.New(0, 0, 0).Sub(inOrig)
	n := vector3.New(0, 1, 0)
	out := reflect(in, n)
	if !out.Equals(vector3.New(-3, 4, -1)) {
		t.Error()
	}
}
