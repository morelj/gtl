release_name := $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD)
release_commit := $(shell git rev-parse --short HEAD)
release_date := $(shell date --iso-8601=seconds)

clean:
	rm -rf build

all: linux-386 linux-amd64 linux-arm linux-arm64 darwin-amd64 windows-386 windows-amd64
	@:

%: export GOOS=$(word 1,$(subst -, ,$@))
%: export GOARCH=$(word 2,$(subst -, ,$@))
%: ext=$(if $(findstring windows,$(GOOS)),.exe)
%:
	@mkdir -p build
	go build -o build/gtl-$(GOOS)-$(GOARCH)$(ext) -ldflags "-X main.ReleaseName=$(release_name) -X main.ReleaseCommit=$(release_commit) -X main.ReleaseDate=$(release_date)"
	cd build && zip gtl-$(GOOS)-$(GOARCH).zip gtl-$(GOOS)-$(GOARCH)$(ext)
