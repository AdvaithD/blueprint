build:
	go build -o bin/blueprint

run: build
	./bin/blueprint

test:
	go test -v ./...
