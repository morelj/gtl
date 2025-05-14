SOURCES = $(shell find -name '*.go') go.mod
DEFAULT_ARCHS = linux-amd64 linux-arm linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

# Find out which is the current architecture
GOOS=$(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH_aarch64_be=arm64be
GOARCH_aarch64=arm64
GOARCH_i386=386
GOARCH_i686=386
GOARCH_x86_64=amd64
GOARCH=$(GOARCH_$(shell uname -m))
ifeq ($(GOARCH),)
	DEFAULT_GOARCH=$(shell uname -m)
endif
EXT_windows=.exe

# Returns the file extension for the target (in the form GOOS-GOARCH)
define get_ext
$(EXT_$(word 1,$(subst -, ,$1)))
endef

BINARIES = $(shell find cmd -mindepth 1 -maxdepth 1 -type d -printf '%P\n')
TARGETS = $(foreach binary,$(BINARIES),$(foreach current,$(DEFAULT_ARCHS),dist/$(binary)-$(current)$(call get_ext,$(current)).zip))

release_name := $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD)
release_commit := $(shell git rev-parse --short HEAD)
release_date := $(shell date --iso-8601=seconds)

default: build/gtl-$(GOOS)-$(GOARCH)$(EXT_$(GOOS))
	@:

all: $(TARGETS)
	@:

build/%: binary=$(word 1,$(subst -, ,$(@F)))
build/%: goos=$(word 2,$(subst -, ,$(@F)))
build/%: goarch=$(word 3,$(subst -, ,$(subst ., ,$(@F))))
build/%: $(SOURCES)
	@mkdir -p $(@D)
	cd cmd/$(binary) && env GOOS=$(goos) GOARCH=$(goarch) CGO_ENABLED=0 go build \
		-ldflags "-X main.ReleaseName=$(release_name) \
			-X main.ReleaseCommit=$(release_commit) \
			-X main.ReleaseDate=$(release_date)" \
		-o ../../$@

dist/%.zip: build/%
	@mkdir -p $(@D)
	cd $(<D) && zip ../$@ $(<F)

clean:
	rm -rf build dist

.PHONY: clean
