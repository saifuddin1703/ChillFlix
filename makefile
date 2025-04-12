# Build commands
build:
	go build -o bin/app cmd/app/main.go

build-prd:
	GOOS=linux GOARCH=amd64 go build -o bin/app cmd/app/main.go

# Run commands
run-local:
	export ENV="local" && ./bin/app

run-prd:
	export ENV="prd" && ./bin/app