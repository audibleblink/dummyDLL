APP = dummy
ARCHES = 386 amd64
ARCH = $(word 1, $@)

cc.386 := i686-w64-mingw32-gcc
cc.amd64 := x86_64-w64-mingw32-gcc

ifeq ($(OS),Windows_NT)
	CC=$(shell go env CC)
else
	CC=${cc.${ARCH}}
endif

all: ${ARCHES}

${ARCHES}:
	@# using go env to set vars since it's OS-agnostic
	@go env -w GOOS=windows
	@go env -w CGO_ENABLED=1
	@go env -w GOARCH=${ARCH}
	@go env -w CC=${CC}

	go build \
		-buildmode=c-shared \
		-trimpath \
		-ldflags="-s -w -buildid= -H windowsgui" \
		-o build/${APP}_${ARCH}.dll

	@# go env makes persistent changes outside the current shell
	@# this undoes those change to preserve a more predicatble
	@go env -u GOOS
	@go env -u CGO_ENABLED
	@go env -u GOARCH
	@go env -u CC

clean:
	@rm -rf build

.PHONY: ${ARCHES} clean
