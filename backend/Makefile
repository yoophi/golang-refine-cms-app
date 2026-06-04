.PHONY: run build tidy test vet clean

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./...

clean:
	rm -rf bin *.db
