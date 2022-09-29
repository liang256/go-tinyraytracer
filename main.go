package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	width, height := 1024, 768
	file, err := os.Create("out.ppm")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(fmt.Sprintf("P6\n%d %d\n255\n", width, height))
	for i := 0; i < height; i++ {
		bi := uint8(float64(i) / float64(height) * 255)
		for j := 0; j < width; j++ {
			bj := uint8(float64(j) / float64(width) * 255)
			file.Write(([]byte{bi, bj, 0})) // RGB
		}
	}
	file.Close()
}
