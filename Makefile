
.PHONY: build all-build clean

ALL_ARCH := amd64 arm64 arm

all: build

build-%:
	@$(MAKE) --no-print-directory ARCH=$* build

all-build: $(addprefix build-, $(ALL_ARCH))

build:
	@scripts/buildapplication
	@scripts/buildcontainer
	@scripts/packageapplication

clean:
	@scripts/cleanup
