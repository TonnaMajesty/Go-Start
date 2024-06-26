package main

import (
	"image"
	"os"
)

func JpegToBgr(filePath string) ([]byte, int, int, error) {
	// 读取JPEG文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 转换为BGR数据
	bgr := make([]byte, width*height*3) // 每个像素RGB各占一个字节

	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 由于上面获取的是16位颜色值，需要转换为8位（1字节）颜色值
			bgr[idx] = byte(b >> 8)
			bgr[idx+1] = byte(g >> 8)
			bgr[idx+2] = byte(r >> 8)
			idx += 3
		}
	}

	return bgr, width, height, nil
}

func JpegToMono8(filePath string) ([]byte, int, int, error) {
	// 读取JPEG文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 转换为Mono8数据
	mono8 := make([]byte, width*height)

	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 取RGB平均值作为灰度值
			gray := uint8((r/0x101 + g/0x101 + b/0x101) / 3)
			mono8[idx] = gray
			idx++
		}
	}

	return mono8, width, height, nil
}
