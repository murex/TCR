all: doc

doc: tidy
	@echo "Generating cobra documentation"
	@go run .

tidy:
	@go mod tidy

.PHONY: all doc tidy