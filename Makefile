build:
	go build -o bin/zip-it ./...

test:
	go test ./...

clean:
	rm -rf bin/

run:
	go run main.go