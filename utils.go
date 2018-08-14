package main

import (
	"fmt"
	"image"
)

func getRawPixelsFromImage(i image.Image) []byte {
	data := make([]byte, 320*640*4)
	width, height := i.Bounds().Max.X, i.Bounds().Max.Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := i.At(x, y).RGBA()
			fmt.Println(r, g, b, a)
			data = append(data, byte(r))
			data = append(data, byte(g))
			data = append(data, byte(b))
			data = append(data, byte(a))
		}
	}
	return data
}
