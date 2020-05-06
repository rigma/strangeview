package main

import (
	"errors"
	"sync"

	"gocv.io/x/gocv"
)

type Facebase struct {
	detector gocv.BRISK
	faces    map[string]faceEntity
	sync.Mutex
}

type faceEntity struct {
	keypoints   []gocv.KeyPoint
	descriptors gocv.Mat
}

func NewFacebase() Facebase {
	return Facebase{
		detector: gocv.NewBRISK(),
		faces:    make(map[string]faceEntity),
	}
}

func (f *Facebase) AddFace(name string, face gocv.Mat) (error, bool) {
	f.Lock()
	defer f.Unlock()

	if _, alreadySaved := f.faces[name]; alreadySaved {
		return errors.New("Face already registered in the facebase!"), false
	}

	keypoints, descriptors := f.detector.DetectAndCompute(face, gocv.NewMat())
	f.faces[name] = faceEntity{
		keypoints:   keypoints,
		descriptors: descriptors,
	}

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
