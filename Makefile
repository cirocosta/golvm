all: install

fmt:
	go fmt
	cd ./lib && go fmt
	cd ./driver && go fmt

install:
	go install -v

test:
	cd ./lib && go test -v

.PHONY: fmt install test
