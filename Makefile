build:
	go build -o bin/godo/godo cmd/godo/main.go

run:
	go run ./cmd/godo/main.go

install:
	go install ./cmd/godo/.
