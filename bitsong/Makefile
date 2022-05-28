export GO111MODULE = on

all: build install

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building sinfonia-bitsong binary..."
	@go build -mod=readonly -o build/sinfonia-bitsong.exe ./cmd/sinfonia-bitsong
else
	@echo "building sinfonia-bitsong binary..."
	@go build -mod=readonly -o build/sinfonia-bitsong ./cmd/sinfonia-bitsong
endif
.PHONY: build

install: go.sum
	@echo "installing sinfonia-bitsong binary..."
	@go install -mod=readonly ./cmd/sinfonia-bitsong
.PHONY: install

clean:
	@rm -rf ./build
.PHONY: clean