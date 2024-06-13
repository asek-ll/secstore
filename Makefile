.PHONY: secstore
secstore:
	go build -o secstore cmd/main.go 

.PHONY: install
install: secstore
	cp ./secstore ~/bin/secstore


.PHONY: clean
clean:
	rm ./secstore
