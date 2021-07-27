package filesystem

type SourceTree interface {
	GetBaseDir() string
	Watch(
		dirList []string,
		filenameMatcher func(filename string) bool,
		interrupt <-chan bool,
	) bool
}
