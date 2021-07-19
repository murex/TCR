package language

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_detect_java_language(t *testing.T) {
	dirPath := filepath.Join("dummy", "java")
	assert.Equal(t, Java{}, DetectLanguage(dirPath))
}

func Test_java_language_name(t *testing.T) {
	assert.Equal(t, "java", Java{}.Name())
}

func Test_list_of_dirs_to_watch_in_java(t *testing.T) {
	var expected = []string{
		filepath.Join("src", "main"),
		filepath.Join("src", "test"),
	}
	assert.Equal(t, expected, DirsToWatch("", Java{}))
}

func Test_filenames_recognized_as_java_src(t *testing.T) {
	expected := []filenameMatching{
		{"Dummy.java", true},
		{"Dummy.JAVA", true},
		{"/dummy/Dummy.java", true},
		{"Dummy.java~", false},
		{"Dummy.java.swp", false},

		{"", false},
		{"dummy", false},
		{"Dummy.cpp", false},
		{"Dummy.sh", false},
	}
	assertFilenames(t, expected, Java{})
}
