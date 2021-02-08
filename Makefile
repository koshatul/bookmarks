
_FILES += $(shell find . -type f -name '*.yml')

README.md: $(_FILES)
	go run main.go
