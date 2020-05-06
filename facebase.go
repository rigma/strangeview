package main

import (
	"errors"
	"sync"

	"gocv.io/x/gocv"
)

type Facebase struct {
	detector gocv.FastFeatureDetector
	faces    map[string][]gocv.KeyPoint
	sync.Mutex
}

func NewFacebase() Facebase {
	return Facebase{
		detector: gocv.NewFastFeatureDetector(),
		faces:    make(map[string][]gocv.KeyPoint),
	}
}

func (f *Facebase) AddFace(name string, face gocv.Mat) (err error, success bool) {
	f.Lock()

	if _, alreadySaved := f.faces[name]; alreadySaved {
		return errors.New("Face already registered in the facebase!"), false
	}
	f.faces[name] = make([]gocv.KeyPoint, 0)

	f.Unlock()

	keypoints := f.detector.Detect(face)

	f.Lock()
	f.faces[name] = keypoints
	f.Unlock()

	return nil, true
}
