package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/hidez8891/shm"
)

func main() {
	Write(0)
	//Read(0)
}

func Write(index int) {
	offset := 1024 * 1024 * 20

	open, _ := os.Open("/Users/tonnamajesty/Downloads/bottle_test_xxxx_front.jpeg")
	defer open.Close()

	img, err := jpeg.Decode(open)
	data := ImageToBgr(img)

	//data, _ := ioutil.ReadAll(open)

	//fmt.Println(len(data))

	r, err := shm.Open("shm_test", 15*1024*1024*30)
	fmt.Println(err)
	defer r.Close()

	n, err := r.WriteAt(data, int64(offset*index))
	fmt.Println(err)
	fmt.Println(n)
}

func Read(index int) {
	offset := 1024 * 1024 * 15

	r, err := shm.Open("shm_test", 15*1024*1024*30)
	fmt.Println(err)
	defer r.Close()

	rbuf := make([]byte, 15*1024*1024*30)
	n, err := r.ReadAt(rbuf, int64(offset*index))

	outfile := fmt.Sprintf("/Users/tonnamajesty/Downloads/bottle_test_xxxx_front.jpeg")

	os.WriteFile(outfile, rbuf, os.ModePerm)

	fmt.Println(err)
	fmt.Println(n)
}

func ImageToBgr(img image.Image) []byte {
	decodedData := make([]byte, img.Bounds().Max.X*img.Bounds().Max.Y*3)

	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y
	pixelIndex := 0
	for h := img.Bounds().Min.Y; h < imgHeight; h++ {
		for w := img.Bounds().Min.X; w < imgWidth; w++ {
			r, g, b, _ := img.At(w, h).RGBA()
			decodedData[pixelIndex+0] = byte(b >> 8)
			decodedData[pixelIndex+1] = byte(g >> 8)
			decodedData[pixelIndex+2] = byte(r >> 8)
			pixelIndex += 3
		}
	}

	return decodedData
}

func ImageToRGB(img image.Image) []byte {
	decodedData := make([]byte, img.Bounds().Max.X*img.Bounds().Max.Y*3)
	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y
	pixelIndex := 0
	for h := 0; h < imgHeight; h++ {
		for w := 0; w < imgWidth; w++ {
			r, g, b, _ := img.At(w, h).RGBA()
			decodedData[pixelIndex+0] = uint8(r >> 8)
			decodedData[pixelIndex+1] = uint8(g >> 8)
			decodedData[pixelIndex+2] = uint8(b >> 8)
			pixelIndex += 3
		}
	}

	return decodedData
}
