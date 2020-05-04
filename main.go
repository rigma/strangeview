package main

import (
	"gocv.io/x/gocv"

	"fmt"
)

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

	for {
		webcam.Read(&input)
		window.IMShow(input)

		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
