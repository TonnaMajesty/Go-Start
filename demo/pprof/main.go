package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io/ioutil"

	"image"
	"image/jpeg"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

func main() {
	go func() {
		for {
			file := "/Users/tonnamajesty/Downloads/1717519070867329761-FIRE_EQUIPMENT-1718908515201860283.jpeg"
			f, _ := os.Open(file)
			data, _ := ioutil.ReadAll(f)
			_, _ = JpegRotate90(data)
			_, _ = JpegRotate902(data)

			time.Sleep(time.Second)
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}

func JpegRotate90(data []byte) ([]byte, error) {
	start := time.Now().UnixMilli()
	img, err := jpeg.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	newImg := Rotate90(img)
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, newImg, nil) // 使用默认的压缩质量，DefaultQuality=75
	fmt.Println("JpegRotate90", time.Now().UnixMilli()-start)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func JpegRotate902(data []byte) ([]byte, error) {
	start := time.Now().UnixMilli()
	img, err := jpeg.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	newImg := Rotate902(img)
	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, newImg, nil) // 使用默认的压缩质量，DefaultQuality=75
	fmt.Println("JpegRotate902", time.Now().UnixMilli()-start)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Rotate902(img image.Image) image.Image {
	return imaging.Rotate(img, 270, color.Black)
}

// 图片旋转90度（顺时针）
func Rotate90(img image.Image) image.Image {
	rotate90 := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dy(), img.Bounds().Dx()))
	for x := img.Bounds().Min.Y; x < img.Bounds().Max.Y; x++ {
		for y := img.Bounds().Max.X - 1; y >= img.Bounds().Min.X; y-- {
			rotate90.Set(img.Bounds().Max.Y-x, y, img.At(y, x))
		}
	}
	return rotate90
}
