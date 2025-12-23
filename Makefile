run:
	go run cmd/api/main.go
build:
	go build -o bin/api cmd/api/main.go
deps:
	go mod download
	go mod tidy

docker-build:
	docker build -t blog-backend .

docker-run:
	docker compose up -d

docker-stop:
	docker compose down

lint:
	go vet ./...

test:
	go test -v ./...
