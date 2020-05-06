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

func (f *Facebase) AddFace(name string, face gocv.Mat) (error, bool) {
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

func (f *Facebase) RemoveFace(name string) (error, bool) {
	f.Lock()
	defer f.Unlock()

	if _, exists := f.faces[name]; !exists {
		return errors.New("Face is not registered in database!"), false
	}

	delete(f.faces, name)

	return nil, true
}

func (f *Facebase) Close() {
	f.detector.Close()
}
