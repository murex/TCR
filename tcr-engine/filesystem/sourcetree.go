package filesystem

// SourceTree is the interface that any implementation must comply with in order to be used
// by TCR engine
type SourceTree interface {
	GetBaseDir() string
	Watch(
		dirList []string,
		filenameMatcher func(filename string) bool,
		interrupt <-chan bool,
	) bool
}
