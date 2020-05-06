package main

import (
	"gocv.io/x/gocv"
)

type Camera struct {
	device *gocv.VideoCapture
	flip   bool
}

func NewCamera() (err error, camera *Camera) {
	device, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		return
	}

	camera = &Camera{
		device: device,
	}
	return
}

func (c *Camera) SetFlip(value bool) {
	c.flip = value
}

func (c *Camera) GetFrame() (frame gocv.Mat) {
	frame = gocv.NewMat()
	c.device.Read(&frame)

	if c.flip {
		tmp := frame.Clone()
		defer tmp.Close()

		gocv.Flip(tmp, &frame, 1)
	}

	return
}
