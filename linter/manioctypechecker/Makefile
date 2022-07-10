.PHONY: mod build-vettool build-golangci-plugin test-vettool test-golangci-plugin all test clean

mod:
	if [ -z "$${GOLANGCI_LINT_TARGET_VERSION}" ]; then \
		GOLANGCI_LINT_TARGET_VERSION=$$(golangci-lint --version | cut -d ' ' -f 4); \
		if [ -z "$${GOLANGCI_LINT_TARGET_VERSION}" ]; then \
			echo "You may need to set the environment variable GOLANGCI_LINT_TARGET_VERSION" >&2; \
			exit 1; \
		fi; \
	fi; \
	( \
		echo "module github.com/fuzmish/manioc/linter/manioctypechecker"; \
		echo ""; \
		echo "go 1.18"; \
		echo ""; \
		echo "require ("; \
		echo "	github.com/golangci/golangci-lint $${GOLANGCI_LINT_TARGET_VERSION}"; \
		echo ")"; \
	) > go.mod
	go mod tidy

build-vettool: mod
	go build -o bin/manioctypechecker cmd/manioctypechecker/main.go

build-golangci-plugin: mod
	go build -buildmode=plugin -o bin/manioctypechecker.so plugin/main.go

test-vettool: build-vettool
	cd testdata/src/a && \
	if [ $$(go vet -vettool=../../../bin/manioctypechecker . 2>&1 \
		| tee /proc/$$$$/fd/2 \
		| grep a.go \
		| wc -l) -ne 55 ]; then exit 1; fi

test-golangci-plugin: build-golangci-plugin
	rm -fr ~/.cache/golangci-lint
	cd testdata/src/a && \
	if [ $$(golangci-lint run --out-format=tab \
		| tee /proc/$$$$/fd/2 \
		| grep manioctypechecker \
		| wc -l) -ne 55 ]; then exit 1; fi

all: build-vettool build-golangci-plugin

test: test-vettool test-golangci-plugin
	go test . -count=1

clean:
	rm -f ./bin/manioctypechecker*
