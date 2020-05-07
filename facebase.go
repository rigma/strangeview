package main

import (
	"errors"
	"runtime"
	"sync"

	"gocv.io/x/gocv"
)

const (
	DETECTION_THRESH  = 40
	LOWE_RATIO_THRESH = .75
)

// A database of "faces" that is storing, in memory, a collection of
// tagged ORB keypoints and descriptors.
type Facebase struct {
	detector gocv.ORB
	matcher  gocv.BFMatcher
	faces    map[string]faceEntity
	sync.Mutex
}

type faceEntity struct {
	keypoints   []gocv.KeyPoint
	descriptors gocv.Mat
}

// Instanciates a new facebase.
func NewFacebase() Facebase {
	return Facebase{
		detector: gocv.NewORB(),
		matcher:  gocv.NewBFMatcher(),
		faces:    make(map[string]faceEntity),
	}
}

// Add a face to the current facebase. You can run this action asynchronously
// to avoid to wait the end of keypoints and descriptors computation.
func (f *Facebase) AddFace(tag string, face gocv.Mat) (error, bool) {
	f.Lock()
	defer f.Unlock()

	if _, alreadySaved := f.faces[tag]; alreadySaved {
		return errors.New("Face already registered in the facebase!"), false
	}

	keypoints, descriptors := f.detector.DetectAndCompute(face, gocv.NewMat())
	f.faces[tag] = faceEntity{
		keypoints:   keypoints,
		descriptors: descriptors,
	}

	return nil, true
}

// Removes a face from the database.
func (f *Facebase) RemoveFace(tag string) (error, bool) {
	f.Lock()
	defer f.Unlock()

	if _, exists := f.faces[tag]; !exists {
		return errors.New("Face is not registered in database!"), false
	}

	delete(f.faces, tag)

	return nil, true
}

// Returns the list of tags stored into the facebase
func (f *Facebase) Tags() (tags []string) {
	for tag := range f.faces {
		tags = append(tags, tag)
	}

	return
}

// Returns a list of detected "face" in an input image
func (f *Facebase) Detect(input gocv.Mat) (err error, faces []Face) {
	_, descriptors := f.detector.DetectAndCompute(input, gocv.NewMat())
	matchesSet := make(map[string][][]gocv.DMatch)

	// We'll try to match the input descriptors with a KNN matching algorithm by deferring
	// the tasks in parallel.
	numCpu := runtime.GOMAXPROCS(0)
	syncChan := make(chan bool, numCpu)

	f.Lock()
	for cpu := 0; cpu < numCpu; cpu++ {
		start, end := cpu*len(f.faces)/numCpu, (cpu+1)*len(f.faces)/numCpu

		go func(faces []string, channel chan bool) {
			for _, face := range faces {
				if _, ok := matchesSet[face]; !ok {
					matchesSet[face] = f.matcher.KnnMatch(descriptors, f.faces[face].descriptors, 2)
				}
			}
			channel <- true
		}(f.Tags()[start:end], syncChan)
	}

	// Waiting for the tasks to be complete
	for cpu := 0; cpu < numCpu; cpu++ {
		<-syncChan
	}
	f.Unlock()

	// Filtering the matches using the Lowe's ratio distance threshold
	for name, matches := range matchesSet {
		var filteredMatches [][]gocv.DMatch

		for _, match := range matches {
			if match[0].Distance < LOWE_RATIO_THRESH*match[1].Distance {
				filteredMatches = append(filteredMatches, match)
			}
		}

		faces = append(faces, Face{
			tag:     name,
			matches: filteredMatches,
		})
	}

	// We'll only indicates that a face is detected when there is
	// enough matches
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

// Close the database
func (f *Facebase) Close() {
	f.detector.Close()
	f.matcher.Close()
}

// A face found into a frame.
type Face struct {
	tag     string
	matches [][]gocv.DMatch
}

// Returns the number of matches in the current face
func (f *Face) MatchesCount() int {
	return len(f.matches)
}

func filter(input []Face, test func(Face) bool) (output []Face) {
	for _, face := range input {
		if test(face) {
			output = append(output, face)
		}
	}

	return
}
