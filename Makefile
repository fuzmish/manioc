.PHONY: coverage lint

coverage:
	if [ ! -d coverage ]; then mkdir coverage; fi
	rm -f coverage/coverage.*
	go test -v -race -coverpkg=./... -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

lint:
	golangci-lint run -v --fix
