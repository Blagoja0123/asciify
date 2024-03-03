package main

import (
	"fmt"
	"image"
	_ "image/png" // Register PNG format decoder
	"math"
	"net/http"
	"os"

	"github.com/Blagoja0123/img-go-ascii/pkg/pixel"
)

func main() {
	// URL of the image
	imageURL := os.Args[1]

	// Send HTTP GET request to fetch the image
	response, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return
	}
	defer response.Body.Close()

	// Decode the response body into an image
	img, _, err := image.Decode(response.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	fmt.Printf("WIDTH: %d; HEIGHT: %d;\n", w, h)
	var pixels [][]pixel.Pixel
	for y := 0; y < h; y++ {
		var pixelRow []pixel.Pixel
		for x := 0; x < w; x++ {
			pixelRow = append(pixelRow, *pixel.NewPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, pixelRow)
	}
	asci := ".,-~:;=!*#$@"

	for _, y := range pixels {
		for _, x := range y {
			char := int(math.Round(toGray(int16(x.Avg()))))
			if char >= 10 {
				char = 9
			}
			// fmt.Printf("int: %d ", char)
			fmt.Printf("%c", asci[char])
		}
		fmt.Printf("\n")
	}
}

func toGray(x int16) float64 {
	return float64(float64(x) / float64(25.5))
}

func compress(pixels [][]pixel.Pixel, factor int) [][]int {
	// Factor stands for an X * X (factor * factor) area of pixels that gets averaged (odd numbers only)
	var result [][]int

	x := factor
	y := x

	xLimit := len(pixels[0]) - factor
	yLimit := len(pixels) - factor
	for ; y < yLimit; y += factor - 2 {
		x = factor
		var row []int
		for ; x < xLimit; x += factor - 2 {
			var avg int
			sum := 0
			//			fmt.Printf("X: %d; Y: %d ", x, y)
			for j := y - factor; j < y+factor-2; j++ {
				for i := x - factor; i < x+factor-2; i++ {
					sum += int(pixels[j][i].Avg())
				}
			}
			avg = sum / (factor * factor)
			row = append(row, avg)
		}
		// fmt.Printf("\n")
		result = append(result, row)
	}

	return result
}
