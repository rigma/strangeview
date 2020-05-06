package main

import (
	"github.com/docker/docker/pkg/namesgenerator"
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
	rawInput := gocv.NewMat()
	facebase := NewFacebase()

	defer rawInput.Close()
	defer facebase.Close()
	defer input.Close()
	defer output.Close()

	for {
		// Reading input from camera and converting into a grayscale image
		webcam.Read(&rawInput)

		input = gocv.NewMat()
		gocv.CvtColor(rawInput, &input, gocv.ColorRGBToGray)

		// Applying a filter on the image
		output = sobelFilter(&input)

		// Displaying the output image
		window.IMShow(output)

		// Checking if a key has been stroked
		key := window.WaitKey(1)
		if key == 27 {
			fmt.Println("Exiting...")
			break
		} else if key == 32 {
			name := namesgenerator.GetRandomName(0)
			fmt.Printf("Saving face %s...\n", name)
			facebase.AddFace(name, rawInput)
		}
	}
}
