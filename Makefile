.PHONY: coverage

coverage:
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -v -coverpkg=./... -coverprofile=coverage/cover.out ./...
	go tool cover -html=coverage/cover.out -o coverage/cover.html
