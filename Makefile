build:
	go build -o bin/main main.go

run:
	go run main.go

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out

lint:
	golangci-lint run --config=.github/linters/golangci.yml

compile:
	GOOS=darwin GOARCH=amd64 go build -o bin/wanikani_cli-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/wanikani_cli-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/wanikani_cli-windows-amd64 main.go
