package engine

import "testing"

func Test_extractTourGuides(t *testing.T) {
	extractTourGuides("./tcr.go", "tour_guide.md", "Driver Round")
}
