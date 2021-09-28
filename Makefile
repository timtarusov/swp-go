all: generate run
generate:
	swag init --parseDependency --parseInternal -g ./cmd/main.go	
run:
	go run ./cmd/.