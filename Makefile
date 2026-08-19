build:
	go build -o ./bin/chainline

run: build
	./bin/chainline
test:
	go test -v ./...

clean:
	rm -rf ./bin/chainline