export GO111MODULE = on

all: build install

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building sinfonia-osmosis binary..."
	@go build -mod=readonly -o build/sinfonia-osmosis.exe ./cmd/sinfonia-osmosis
else
	@echo "building sinfonia-osmosis binary..."
	@go build -mod=readonly -o build/sinfonia-osmosis ./cmd/sinfonia-osmosis
endif
.PHONY: build

install: go.sum
	@echo "installing sinfonia-osmosis binary..."
	@go install -mod=readonly ./cmd/sinfonia-osmosis
.PHONY: install

clean:
	@rm -rf ./build
.PHONY: clean