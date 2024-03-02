package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg" // Register JPEG format decoder
	_ "image/png"  // Register PNG format decoder
	"math"
	"net/http"
	"os"
)

func main() {
	// URL of the image
	imageURL := "https://static.scientificamerican.com/sciam/cache/file/2AE14CDD-1265-470C-9B15F49024186C10_source.jpg?w=600"

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
	var pixels [][]Pixel
	for y := 0; y < h; y++ {
		var pixelRow []Pixel
		for x := 0; x < w; x++ {
			pixelRow = append(pixelRow, *NewPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, pixelRow)
	}
	asci := ".,-~:;=!*#$@"
	for _, y := range pixels {
		for _, x := range y {
			char := int(math.Round(toGray(x.avg())))
			fmt.Printf("%c", asci[char])
		}
		fmt.Printf("\n")
	}
	// Create a new file to save the image
	file, err := os.Create("downloaded_image.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Encode the image and save it to the file
	err = jpeg.Encode(file, img, nil)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Image downloaded successfully!")
}

type Pixel struct {
	R, G, B, A uint32
}

func NewPixel(r, g, b, a uint32) *Pixel {
	return &Pixel{
		R: r, G: g, B: b, A: a,
	}
}

func (p *Pixel) avg() int16 {
	return (int16(p.R/255) + int16(p.G/255) + int16(p.B/255)) / 3
}

func toGray(x int16) float64 {
	return float64(float64(x) / float64(25.5))
}
