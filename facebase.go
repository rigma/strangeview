package main

import (
	"errors"
	"sync"

	"gocv.io/x/gocv"
)

const (
	DETECTION_THRESH  = 40
	LOWE_RATIO_THRESH = .75
)

type Facebase struct {
	detector gocv.BRISK
	matcher  gocv.BFMatcher
	faces    map[string]faceEntity
	sync.Mutex
}

func NewFacebase() Facebase {
	return Facebase{
		detector: gocv.NewBRISK(),
		matcher:  gocv.NewBFMatcher(),
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

func (f *Facebase) Detect(input gocv.Mat) (err error, faces []Face) {
	_, descriptors := f.detector.DetectAndCompute(input, gocv.NewMat())
	matchesSet := make(map[string][][]gocv.DMatch)

	f.Lock()
	for name, face := range f.faces {
		matchesSet[name] = f.matcher.KnnMatch(descriptors, face.descriptors, 2)
	}
	f.Unlock()

	for name, matches := range matchesSet {
		var filteredMatches [][]gocv.DMatch

		for _, match := range matches {
			if match[0].Distance < LOWE_RATIO_THRESH*match[1].Distance {
				filteredMatches = append(filteredMatches, match)
			}
		}

		faces = append(faces, Face{
			name:    name,
			matches: filteredMatches,
		})
	}

	faces = filter(faces, func(face Face) bool {
		return face.MatchesCount() >= DETECTION_THRESH
	})

	if len(faces) == 0 {
		err = errors.New("No faces are found!")
		faces = nil

		return
	}

	err = nil
	return
}

func (f *Facebase) Close() {
	f.detector.Close()
	f.matcher.Close()
}

type Face struct {
	name    string
	matches [][]gocv.DMatch
}

func (f *Face) MatchesCount() int {
	return len(f.matches)
}

type faceEntity struct {
	keypoints   []gocv.KeyPoint
	descriptors gocv.Mat
}

func filter(input []Face, test func(Face) bool) (output []Face) {
	for _, face := range input {
		if test(face) {
			output = append(output, face)
		}
	}

	return
}