# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gzen android ios gzen-cross evm all test clean docs
.PHONY: gzen-linux gzen-linux-386 gzen-linux-amd64 gzen-linux-mips64 gzen-linux-mips64le
.PHONY: gzen-linux-arm gzen-linux-arm-5 gzen-linux-arm-6 gzen-linux-arm-7 gzen-linux-arm64
.PHONY: gzen-darwin gzen-darwin-386 gzen-darwin-amd64
.PHONY: gzen-windows gzen-windows-386 gzen-windows-amd64
.PHONY: gzen all test lint fmt clean devtools help

GO ?= latest
GOBIN = $(CURDIR)/build/bin
GORUN = env GO111MODULE=on go run
GOPATH = $(shell go env GOPATH)

GIT_COMMIT ?= $(shell git rev-list -1 HEAD)

PACKAGE = github.com/zenanet-network/go-zenanet
GO_FLAGS += -buildvcs=false
GO_LDFLAGS += -ldflags "-X ${PACKAGE}/params.GitCommit=${GIT_COMMIT}"

TESTALL = $$(go list ./... | grep -v go-zenanet/cmd/)
TESTE2E = ./tests/...
GOTEST = GODEBUG=cgocheck=0 go test $(GO_FLAGS) $(GO_LDFLAGS) -p 1

zena:
	mkdir -p $(GOPATH)/bin/
	go build -o $(GOBIN)/zena $(GO_LDFLAGS) ./cmd/cli/main.go
	cp $(GOBIN)/zena $(GOPATH)/bin/
	@echo "Done building."

protoc:
	protoc --go_out=. --go-grpc_out=. ./internal/cli/server/proto/*.proto

generate-mocks:
	go generate mockgen -destination=./tests/zena/mocks/IHarmoniaClient.go -package=mocks ./consensus/zena IHarmoniaClient
	go generate mockgen -destination=./eth/filters/IDatabase.go -package=filters ./ethdb Database
	go generate mockgen -destination=./eth/filters/IBackend.go -package=filters ./eth/filters Backend
	go generate mockgen -destination=../eth/filters/IDatabase.go -package=filters ./ethdb Database

#? gzen: Build gzen.
gzen:
	$(GORUN) build/ci.go install ./cmd/gzen
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gzen\" to launch gzen."

#? all: Build all packages and executables.
all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gzen.aar\" to use the library."
	@echo "Import \"$(GOBIN)/gzen-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"

ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gzen.framework\" to use the library."

test:
	$(GOTEST) --timeout 30m -cover -short -coverprofile=cover.out -covermode=atomic $(TESTALL)

test-txpool-race:
	$(GOTEST) -run=TestPoolMiningDataRaces --timeout 600m -race -v ./core/

test-race:
	$(GOTEST) --timeout 15m -race -shuffle=on $(TESTALL)

gocovmerge-deps:
	$(GOBUILD) -o $(GOBIN)/gocovmerge github.com/wadey/gocovmerge

test-integration:
	$(GOTEST) --timeout 60m -cover -coverprofile=cover.out -covermode=atomic -tags integration $(TESTE2E)

escape:
	cd $(path) && go test -gcflags "-m -m" -run none -bench=BenchmarkJumpdest* -benchmem -memprofile mem.out

lint:
	@./build/bin/golangci-lint run --config ./.golangci.yml

lintci-deps:
	rm -f ./build/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./build/bin v1.57.2

goimports:
	goimports -local "$(PACKAGE)" -w .

docs:
	$(GORUN) cmd/clidoc/main.go -d ./docs/cli

#? fmt: Ensure consistent code formatting.
fmt:
	gofmt -s -w $(shell find . -name "*.go")

#? clean: Clean go cache, built executables, and the auto generated folder.
clean:
	go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

#? devtools: Install recommended developer tools.
devtools:
	# Notice! If you adding new binary - add it also to tests/deps/fake.go file
	$(GOBUILD) -o $(GOBIN)/stringer github.com/golang.org/x/tools/cmd/stringer
	$(GOBUILD) -o $(GOBIN)/go-bindata github.com/kevinburke/go-bindata/go-bindata
	$(GOBUILD) -o $(GOBIN)/codecgen github.com/ugorji/go/codec/codecgen
	$(GOBUILD) -o $(GOBIN)/abigen ./cmd/abigen
	$(GOBUILD) -o $(GOBIN)/mockgen github.com/golang/mock/mockgen
	$(GOBUILD) -o $(GOBIN)/protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go
	PATH=$(GOBIN):$(PATH) go generate ./common
	PATH=$(GOBIN):$(PATH) go generate ./core/types
	PATH=$(GOBIN):$(PATH) go generate ./consensus/eirene
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)
gzen-cross: gzen-linux gzen-darwin gzen-windows gzen-android gzen-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gzen-*

gzen-linux: gzen-linux-386 gzen-linux-amd64 gzen-linux-arm gzen-linux-mips64 gzen-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-*

gzen-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gzen
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep 386

gzen-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gzen
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep amd64

gzen-linux-arm: gzen-linux-arm-5 gzen-linux-arm-6 gzen-linux-arm-7 gzen-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep arm

gzen-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gzen
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep arm-5

gzen-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gzen
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep arm-6

gzen-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gzen
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep arm-7

gzen-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gzen
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep arm64

gzen-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gzen
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep mips

gzen-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gzen
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep mipsle

gzen-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gzen
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep mips64

gzen-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gzen
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gzen-linux-* | grep mips64le

gzen-darwin: gzen-darwin-386 gzen-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gzen-darwin-*

gzen-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gzen
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-darwin-* | grep 386

gzen-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gzen
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-darwin-* | grep amd64

gzen-windows: gzen-windows-386 gzen-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gzen-windows-*

gzen-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gzen
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-windows-* | grep 386

gzen-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gzen
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gzen-windows-* | grep amd64

PACKAGE_NAME          := github.com/maticnetwork/bor
GOLANG_CROSS_VERSION  ?= v1.22.1

.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e GITHUB_TOKEN \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate --skip-publish

.PHONY: release
release:
	@docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e GITHUB_TOKEN \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-e SLACK_WEBHOOK \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(HOME)/.docker/config.json:/root/.docker/config.json \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate
