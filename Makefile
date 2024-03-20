.PHONY: build
build:
	go build -o secstore cmd/main.go 

install:
	cp ./secstore ~/bin/secstore
