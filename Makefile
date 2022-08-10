appname := homescript
workingdir := "homescript-cli"
sources := $(wildcard *.go)
version := 2.14.0

build = CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -v -ldflags "-s -w" -o $(appname)$(3) $(4)
tar = mkdir -p build && tar -cvzf ./$(appname)_$(1)_$(2).tar.gz $(appname)$(3) && mv $(appname)_$(1)_$(2).tar.gz build

# Clean
clean:
	rm -rf bin
	rm -rf build
	rm -f homescript

# Change version
version:
	python3 update_version.py

# Github release
gh-release:
	gh release create v$(version) ./build/*.tar.gz -F ./CHANGELOG.md -t 'CLI v$(version)'

lint:
	go vet
	golangci-lint run
	typos

# Release (currently only build)
release: clean lint build

# Builds
build: clean linux

# Build architectures
linux: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(sources)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64, -installsuffix cgo)
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

