package main

import (
	"gocv.io/x/gocv"

	"fmt"
	"image"
)

func sobelFilter(input *gocv.Mat) (output gocv.Mat) {
	blured_input := gocv.NewMat()
	defer blured_input.Close()

	gocv.GaussianBlur(*input, &blured_input, image.Pt(3, 3), 0, 0, gocv.BorderDefault)

	grad_x, grad_y := gocv.NewMat(), gocv.NewMat()
	defer grad_x.Close()
	defer grad_y.Close()

	gocv.Sobel(blured_input, &grad_x, -1, 1, 0, 3, 1.0, 0.0, gocv.BorderDefault)
	gocv.Sobel(blured_input, &grad_y, -1, 0, 1, 3, 1.0, 0.0, gocv.BorderDefault)

	output = gocv.NewMat()
	gocv.AddWeighted(grad_x, 0.5, grad_y, 0.5, 0, &output)
	return
}

func main() {
	fmt.Println("Launching GoView...")

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		panic("Unable to retrieve the video capture device!")
	}

	window := gocv.NewWindow("GoView")
	defer window.Close()

	input := gocv.NewMat()
	defer input.Close()

	output := gocv.NewMat()
	defer output.Close()

	for {
		webcam.Read(&input)
		output = sobelFilter(&input)
		window.IMShow(output)

		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
