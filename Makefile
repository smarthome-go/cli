appname := homescript
workingdir := "homescript-cli"
sources := $(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -o $(appname)$(3) $(4)
tar = mkdir -p build && tar -cvzf ./$(appname)_$(1)_$(2).tar.gz $(appname)$(3) && mv $(appname)_$(1)_$(2).tar.gz build

# Clean
clean:
	rm -rf bin
	rm -rf build

# Builds
build: clean linux

# Build architectures
linux: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(sources)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64, -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build')
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

