init:
	git submodule init
	git submodule update

run:
	go run .

build :
	go build -o verifier .
