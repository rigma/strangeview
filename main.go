package main

import (
	"github.com/docker/docker/pkg/namesgenerator"
	"gocv.io/x/gocv"

	"fmt"
	"image"
)

// Returns a blured image matrix thanks to a Gaussian blur filter.
func blur(input *gocv.Mat) (output gocv.Mat) {
	output = gocv.NewMat()

	gocv.GaussianBlur(*input, &output, image.Pt(3, 3), 1.0, 1.0, gocv.BorderDefault)
	return
}

// Returns a filtered matrix thanks to a Sobel filter. Therefore, the edges of the
// input are extracted on the output image.
func sobelFilter(input *gocv.Mat) (output gocv.Mat) {
	// Then we'll apply the Sobel filter following the X and the Y direction
	gradX, gradY := gocv.NewMat(), gocv.NewMat()
	defer gradX.Close()
	defer gradY.Close()

	gocv.Sobel(*input, &gradX, -1, 1, 0, 3, 1.0, 0.0, gocv.BorderDefault)
	gocv.Sobel(*input, &gradY, -1, 0, 1, 3, 1.0, 0.0, gocv.BorderDefault)

	// Finally we'll do a mean of the two matrices and return it
	output = gocv.NewMat()
	gocv.AddWeighted(gradX, 0.5, gradY, 0.5, 0, &output)
	return
}

func main() {
	fmt.Println("Launching StrangeView...")

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		panic("Unable to retrieve the video capture device!")
	}

	window := gocv.NewWindow("StrangeView")
	defer window.Close()

	rawInput := gocv.NewMat()
	facebase := NewFacebase()
	sobel := false

	defer rawInput.Close()
	defer facebase.Close()

	for {
		// Reading input from camera and converting into a grayscale image
		webcam.Read(&rawInput)
		rawInput = blur(&rawInput)

		if sobel {
			filtered := gocv.NewMat()
			gocv.CvtColor(rawInput, &filtered, gocv.ColorRGBToGray)
			filtered = sobelFilter(&filtered)

			window.IMShow(filtered)
		} else {
			window.IMShow(rawInput)
		}

		// Checking if a key has been stroked
		key := window.WaitKey(1)
		if key == 27 {
			fmt.Println("Exiting...")
			break
		} else if key == 32 {
			name := namesgenerator.GetRandomName(0)
			fmt.Printf("Saving face %s...\n", name)
			facebase.AddFace(name, rawInput)
		} else if key == 13 {
			err, faces := facebase.Detect(rawInput)
			if err != nil {
				fmt.Println(err)
			} else {
				for _, face := range faces {
					fmt.Printf("Face %s detected!\n", face.name)
				}
			}
		} else if key == 49 {
			fmt.Println("Toggling Sobel filter...")
			sobel = !sobel
		}
	}
}
