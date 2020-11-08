build:
	go build -o bin/main main.go

run:
	go run main.go

test:
	go test -v ./... -covermode=count -coverprofile=coverage.out

lint:
	golangci-lint run --config=.github/linters/golangci.yml

clean:
	rm -r bin/**

compile:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/wanikani_cli main.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/wanikani_cli main.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/wanikani_cli.exe main.go
