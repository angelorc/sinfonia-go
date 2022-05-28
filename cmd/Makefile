export GO111MODULE = on

all: build install

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building sinfonia binary..."
	@go build -mod=readonly -o build/sinfonia.exe ./sinfonia
else
	@echo "building sinfonia binary..."
	@go build -mod=readonly -o build/sinfonia ./sinfonia
endif
.PHONY: build

install: go.sum
	@echo "installing sinfonia binary..."
	@go install -mod=readonly ./sinfonia
.PHONY: install

clean:
	@rm -rf ./build
.PHONY: clean