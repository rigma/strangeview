package main

import (
	"gocv.io/x/gocv"

	"fmt"
	"image"
)

// Returns a filtered matrix thanks to a Sobel filter. Therefore, the edges of the
// input are extracted on the output image.
func sobelFilter(input *gocv.Mat) (output gocv.Mat) {
	// First, we'll blur the input image to remove signal noise of the camera
	blured_input := gocv.NewMat()
	defer blured_input.Close()

	gocv.GaussianBlur(*input, &blured_input, image.Pt(3, 3), 1.0, 0, gocv.BorderDefault)

	// Then we'll apply the Sobel filter following the X and the Y direction
	grad_x, grad_y := gocv.NewMat(), gocv.NewMat()
	defer grad_x.Close()
	defer grad_y.Close()

	gocv.Sobel(blured_input, &grad_x, -1, 1, 0, 3, 1.0, 0.0, gocv.BorderDefault)
	gocv.Sobel(blured_input, &grad_y, -1, 0, 1, 3, 1.0, 0.0, gocv.BorderDefault)

	// Finally we'll do a mean of the two matrices and return it
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

	var input, output gocv.Mat
	raw_input := gocv.NewMat()

	defer raw_input.Close()
	defer input.Close()
	defer output.Close()

	for {
		// Reading input from camera and converting into a grayscale image
		webcam.Read(&raw_input)

		input = gocv.NewMat()
		gocv.CvtColor(raw_input, &input, gocv.ColorRGBToGray)

		// Applying a filter on the image
		output = sobelFilter(&input)

		// Displaying the output image
		window.IMShow(output)

		// Checking if a key has been stroked
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
