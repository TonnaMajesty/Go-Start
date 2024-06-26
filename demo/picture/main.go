package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	// test jpeg to mono8
	mono8filePath := "/Users/tonnamajesty/Documents/酒瓶测试素材/1700721819447-YH0210110003.jpeg"
	mono8Data, width, height, err := JpegToMono8(mono8filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return

	}
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)

	imgMono8, err := Mono8ToJpegParallel(mono8Data, width, height)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 保存img

	mono8f, err := os.Create("/Users/tonnamajesty/Documents/酒瓶测试素材/1700721819447-YH0210110003-zhuanhuan.jpeg")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer mono8f.Close()
	err = jpeg.Encode(mono8f, imgMono8, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return

	// test jpeg to bgr
	filePath := "/Users/tonnamajesty/Documents/酒瓶测试素材/20240506-112354.jpeg"
	bgrData, width, height, err := JpegToBgr(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Width:", width)
	fmt.Println("Height:", height)
	imgBgr, err := BGRToJpegParallel(bgrData, width, height)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 保存img

	f, err := os.Create("/Users/tonnamajesty/Documents/酒瓶测试素材/20240506-112354-zhuanhuan.jpg")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()
	err = jpeg.Encode(f, imgBgr, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return

	//open1, _ := os.Open("/Users/tonnamajesty/Downloads/20231027-134242.jpeg")
	//defer open1.Close()
	//b, err := JpgToMono8(open1)
	//
	//outfile1 := fmt.Sprintf("/Users/tonnamajesty/Downloads/mz_test_%d.bin", time.Now().UnixMilli())
	//os.WriteFile(outfile1, b, os.ModePerm)

	open1, _ := os.Open("/Users/tonnamajesty/Downloads/mz_test_1698392725228.bin")
	defer open1.Close()

	b, err := io.ReadAll(open1)

	bu, err := Mono8ToJpegParallelWithCrop(b, 2448, 2048, 0, 0, 2448, 2048)

	outfile1 := fmt.Sprintf("/Users/tonnamajesty/Downloads/mz_test_%d.jpeg", time.Now().UnixMilli())
	os.WriteFile(outfile1, bu.Bytes(), os.ModePerm)

	return
	open, _ := os.Open("/Users/tonnamajesty/Downloads/bottle_front.jpg")
	//open, _ := os.Open("/Users/tonnamajesty/Downloads/Pic_2023_09_04_161425_32.jpeg")
	//open, _ := os.Open("/Users/tonnamajesty/Downloads/raw_test_1679546767080")

	defer open.Close()

	//data, _ := ioutil.ReadAll(open)
	//
	//for i := 0; i < 1; i++ {
	//	now := time.Now().UnixMilli()
	//	file, _ := BGRToJpeg(data, 2448, 2048)
	//	fmt.Println(time.Now().UnixMilli() - now)
	//	outfile := fmt.Sprintf("/Users/tonnamajesty/Downloads/raw_test_%d.jpeg", time.Now().UnixMilli())
	//	os.WriteFile(outfile, file.Bytes(), os.ModePerm)
	//}

	img, _ := jpeg.Decode(open)
	data := ImageToBgr(img)
	jpg, err := BGRToJpeg(data, 2448, 2048)
	//jpg, err := BGRToJpeg(data, 12000, 4096)
	fmt.Println(err)

	outfile := fmt.Sprintf("/Users/tonnamajesty/Downloads/bottle_test_%d.jpeg", time.Now().UnixMilli())
	//rawfile := fmt.Sprintf("/Users/tonnamajesty/Downloads/raw_test_%d", time.Now().UnixMilli())

	os.WriteFile(outfile, jpg.Bytes(), os.ModePerm)
	//os.WriteFile(rawfile, data, os.ModePerm)
}

func Crop() {
	var err error
	//open, _ := os.Open("/Users/tonnamajesty/Documents/酒瓶测试素材/test_xs_1695711381144.jpeg")
	//defer open.Close()

	file, _ := os.Open("/Users/tonnamajesty/Documents/酒瓶测试素材/test_rotate_1695697649966.bin")
	defer file.Close()
	fileBytes, _ := io.ReadAll(file)

	//img, err := imaging.Decode(open)
	//fmt.Println(err)
	//

	buffer := &bytes.Buffer{}

	//img, err := BGRToJpegParallelWithCrop(fileBytes, 12000, 4096, 7208, 175, 2015, 2579)
	//fmt.Println(err)

	img := BgrToJpegWithCrop(fileBytes, 12000, 4096, 7208, 175, 2015, 2579)

	//jpeg.Encode(buffer, img, &jpeg.Options{Quality: 75})

	//img, _ = imaging.Decode(buffer)
	//
	//img = imaging.Crop(img, image.Rect(7208, 175, 9223, 2754))

	dc := gg.NewContextForImage(img)
	//dc.SetRGB(255, 255, 240)            // Set the line color to red
	dc.SetLineWidth(10) // Set the line width to 2 pixels
	dc.SetColor(color.NRGBA{R: 255, G: 255, B: 0, A: 255})
	dc.DrawRectangle(1612, 0, 315, 165) // Replace x, y, width, and height with the desired values
	dc.Stroke()
	img = dc.Image()

	err = jpeg.Encode(buffer, img, &jpeg.Options{Quality: 75})
	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile(fmt.Sprintf("/Users/tonnamajesty/Downloads/test_crop_%d.jpeg", time.Now().UnixMilli()), buffer.Bytes(), os.ModePerm)
}

func BGRToJpeg(content []byte, width, height int) (*bytes.Buffer, error) {
	if width == 0 || height == 0 {
		width = 2448
		height = 2048
	}
	buffer := &bytes.Buffer{}
	//if len(content) != width*height*3 {
	//	logrus.Errorf("BGRToJpeg,content.length:%v!= width:%v*height:%v*3", len(content), width, height)
	//	return buffer, fmt.Errorf("invalid content length, expect:%d, got:%d", width*height*3, len(content))
	//}

	now := time.Now().UnixMilli()
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	x, y := 0, 0
	for i := 0; i < len(content); i += 3 {
		b, g, r := content[i], content[i+1], content[i+2]
		rgba.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 0})
		if x > 0 && x%(width-1) == 0 {
			x = 0
			y++
		} else {
			x++
		}
	}

	fmt.Printf("for %d\n", time.Now().UnixMilli()-now)

	err := jpeg.Encode(buffer, rgba, &jpeg.Options{Quality: 100})
	if err != nil {
		logrus.Errorf("jpeg.Encode,err:%+v", err)
		return buffer, err
	}

	return buffer, nil
}

func BGRToJpegV2(content []byte, width, height int) (*bytes.Buffer, error) {
	if width == 0 || height == 0 {
		width = 2488
		height = 2044
	}
	buffer := &bytes.Buffer{}
	//if len(content) != width*height*3 {
	//	logrus.Errorf("BGRToJpeg,content.length:%v!= width:%v*height:%v*3", len(content), width, height)
	//	return buffer, fmt.Errorf("invalid content length, expect:%d, got:%d", width*height*3, len(content))
	//}

	now := time.Now().UnixMilli()
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	x, y := 0, 0

	var wg sync.WaitGroup
	numParts := 10
	partSize := len(content) / numParts
	for i := 0; i < numParts; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := start; j < start+partSize; j += 3 {
				if j > 0 && j%3 == 0 {
					b, g, r := content[j-3], content[j-2], content[j-1]
					rgba.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 0})
					if x > 0 && x%(width-1) == 0 {
						x = 0
						y++
					} else {
						x++
					}
				}
			}
		}(i * partSize)
	}
	wg.Wait()

	fmt.Printf("for %d\n", time.Now().UnixMilli()-now)

	err := jpeg.Encode(buffer, rgba, &jpeg.Options{Quality: 100})
	if err != nil {
		logrus.Errorf("jpeg.Encode,err:%+v", err)
		return buffer, err
	}

	return buffer, nil
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

func BGRToJpegParallel(content []byte, width, height int) (image.Image, error) { // W: exported function BGRToJpeg should have comment or be unexported
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	routineNum := 4
	step := width * height / routineNum * 3
	for i := 0; i < routineNum; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			rgbaWorker(content, rgba, i, step, width)
		}()
	}
	wg.Wait()

	return rgba, nil
}

func rgbaWorker(content []byte, rgba *image.RGBA, workerIndex int, workerStep int, imgWidth int) {
	contentLen := len(content)
	i := workerIndex * workerStep
	alpha := uint8(0)
	stride := imgWidth
	for ; i < (workerIndex+1)*workerStep && i < contentLen; i += 12 {
		b0 := content[i+0]
		g0 := content[i+1]
		r0 := content[i+2]
		b1 := content[i+3]
		g1 := content[i+4]
		r1 := content[i+5]
		b2 := content[i+6]
		g2 := content[i+7]
		r2 := content[i+8]
		b3 := content[i+9]
		g3 := content[i+10]
		r3 := content[i+11]

		xy := i / 3
		rgba.Set((xy+0)%stride, (xy+0)/stride, color.RGBA{R: r0, G: g0, B: b0, A: alpha})
		rgba.Set((xy+1)%stride, (xy+1)/stride, color.RGBA{R: r1, G: g1, B: b1, A: alpha})
		rgba.Set((xy+2)%stride, (xy+2)/stride, color.RGBA{R: r2, G: g2, B: b2, A: alpha})
		rgba.Set((xy+3)%stride, (xy+3)/stride, color.RGBA{R: r3, G: g3, B: b3, A: alpha})
	}
	for ; i < (workerIndex+1)*workerStep && i < contentLen; i += 3 {
		b0 := content[i+0]
		g0 := content[i+1]
		r0 := content[i+2]

		xy := i / 3
		rgba.Set((xy+0)%stride, (xy+0)/stride, color.RGBA{R: r0, G: g0, B: b0, A: alpha})
	}
}

func BgrToJpegWithCrop(content []byte, width, height int, cropX, cropY, cropWidth, cropHeight int) image.Image {
	// Create a new RGBA image
	rgba := image.NewRGBA(image.Rect(0, 0, cropWidth, cropHeight))

	// Set the number of goroutines to use for parallel processing
	numGoroutines := 4

	// Calculate the number of pixels each goroutine will process
	pixelsPerGoroutine := (cropWidth * cropHeight) / numGoroutines

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Process each region of pixels in parallel
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startIndex, endIndex int) {
			defer wg.Done()
			processPixels(startIndex, endIndex, width, height, cropX, cropY, cropWidth, cropHeight, content, rgba)
		}(i*pixelsPerGoroutine, (i+1)*pixelsPerGoroutine)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return rgba
}

func processPixels(startIndex, endIndex, width, height, x, y, cropWidth, cropHeight int, bgrData []byte, rgba *image.RGBA) {
	for i := startIndex; i < endIndex; i++ {
		// Calculate the x and y coordinates of the current pixel
		pixelX := i % cropWidth
		pixelY := i / cropWidth

		// Calculate the actual coordinates in the original image
		imageX := x + pixelX
		imageY := y + pixelY

		// Check if the current pixel falls within the crop region
		if imageX >= 0 && imageX < width && imageY >= 0 && imageY < height {
			// Calculate the index of the current pixel in the BGR byte data
			index := (imageY*width + imageX) * 3

			// Get the B, G, R values from the BGR byte data
			b := bgrData[index]
			g := bgrData[index+1]
			r := bgrData[index+2]

			// Set the pixel value in the RGBA image
			rgba.Set(pixelX, pixelY, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
}

func JpgToMono8(content io.Reader) ([]byte, error) {
	// 解码JPEG图像
	img, err := jpeg.Decode(content)
	if err != nil {
		return nil, err
	}

	// 获取图像的尺寸
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	fmt.Println(width, height)

	// 创建一个字节切片来存储Mono8图像数据
	mono8Data := make([]byte, width*height)

	// 遍历图像的每个像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取当前像素的颜色
			grayColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			// 将灰度值存储到Mono8数据中
			mono8Data[y*width+x] = grayColor.Y
		}
	}

	// 创建一个新的Mono8图像
	mono8 := image.NewGray(img.Bounds())

	// 将每个像素从RGB转换为灰度值
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// 获取RGB像素值
			r, g, b, _ := img.At(x, y).RGBA()

			// 计算灰度值
			gray := uint8((r + g + b) / 3 >> 8)

			// 设置Mono8图像的像素值
			mono8.SetGray(x, y, color.Gray{gray})
		}
	}

	return mono8Data, nil
}

func Mono8ToJpegParallelWithCrop(content []byte, imgWidth, imgHeight, cropX, cropY, cropWidth, cropHeight int) (bytes.Buffer, error) {
	buffer := bytes.Buffer{}
	img := image.NewGray(image.Rect(0, 0, cropWidth, cropHeight))

	semaphore := make(chan struct{}, 4)

	var wg sync.WaitGroup

	for y := cropY; y < cropY+cropHeight; y++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(y int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			for x := cropX; x < cropX+cropWidth; x++ {
				value := content[y*imgWidth+x]
				gray := color.Gray{Y: value}
				img.Set(x-cropX, y-cropY, gray)
			}
		}(y)
	}

	wg.Wait()

	err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return buffer, errors.Wrap(err, "Mono8ToJpegParallel")
	}

	return buffer, nil
}

func Mono8ToJpegParallel(content []byte, imgWidth, imgHeight int) (image.Image, error) {
	//buffer := bytes.Buffer{}
	img := image.NewGray(image.Rect(0, 0, imgWidth, imgHeight))

	semaphore := make(chan struct{}, 4)

	var wg sync.WaitGroup

	for y := 0; y < imgHeight; y++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(y int) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			for x := 0; x < imgWidth; x++ {
				value := content[y*imgWidth+x]
				gray := color.Gray{Y: value}
				img.Set(x, y, gray)
			}
		}(y)
	}

	wg.Wait()

	//err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: 100})
	//if err != nil {
	//	return buffer, errors.Wrap(err, "Mono8ToJpegParallel")
	//}
	return img, nil
}
