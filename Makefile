all: build


dockerbuild:
	cd provider && go mod tidy && go build && 	swag init &&	docker build --platform linux/amd64 -f DOCKERFILE . -t provider
	cd server &&  go mod tidy && go build  &&	swag init && 	docker build --platform linux/amd64 -f DOCKERFILE . -t server
	docker compose -f docker-compose.yaml up --build -d