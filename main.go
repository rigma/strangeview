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
	// First we'll convert the raw input into a grayscale image
	grayscale := gocv.NewMat()
	defer grayscale.Close()

	gocv.CvtColor(*input, &grayscale, gocv.ColorRGBToGray)

	// Then we'll apply the Sobel filter following the X and the Y direction
	gradX, gradY := gocv.NewMat(), gocv.NewMat()
	defer gradX.Close()
	defer gradY.Close()

	// Sobel filter in X direction
	gocv.Sobel(grayscale, &gradX, -1, 1, 0, 3, 1.0, 0.0, gocv.BorderDefault)
	// Sobel filter in Y direction
	gocv.Sobel(grayscale, &gradY, -1, 0, 1, 3, 1.0, 0.0, gocv.BorderDefault)

	// Finally we'll do a mean of the two matrices and return it
	output = gocv.NewMat()
	gocv.AddWeighted(gradX, 0.5, gradY, 0.5, 0, &output)
	return
}

// Returns an inverted Sobel filter of the input image which looks like
// a drawing on a canvas.
func canvasFilter(input *gocv.Mat) (output gocv.Mat) {
	output = sobelFilter(input)
	gocv.Threshold(output, &output, 18, 255, gocv.ThresholdBinaryInv)
	return
}

func main() {
	fmt.Println("Launching StrangeView...")

	err, camera := NewCamera()
	if err != nil {
		panic("Unable to retrieve the video capture device!")
	}

	camera.SetFlip(true)
	window := gocv.NewWindow("StrangeView")
	defer window.Close()

	var sobel, canvas bool
	var rawInput gocv.Mat
	facebase := NewFacebase()

	defer rawInput.Close()
	defer facebase.Close()

	for {
		// Reading input from camera and converting into a grayscale image
		rawInput = camera.GetFrame()
		rawInput = blur(&rawInput)

		if sobel {
			output := sobelFilter(&rawInput)
			window.IMShow(output)
			output.Close()
		} else if canvas {
			output := canvasFilter(&rawInput)
			window.IMShow(output)
			output.Close()
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
			go facebase.AddFace(name, rawInput)
		} else if key == 13 {
			err, faces := facebase.Detect(rawInput)
			if err != nil {
				fmt.Println(err)
			} else {
				for _, face := range faces {
					fmt.Printf("Face %s detected!\n", face.tag)
				}
			}
		} else if key == 49 && (sobel || canvas) {
			fmt.Println("Normal view...")
			sobel, canvas = false, false
		} else if key == 50 && !sobel {
			fmt.Println("Sobel filtered view")
			sobel, canvas = true, false
		} else if key == 51 && !canvas {
			fmt.Println("Canvas filtered view")
			sobel, canvas = false, true
		}
	}
}
