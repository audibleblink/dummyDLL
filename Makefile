APP=dummy
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
	go env -w GOOS=windows
	go env -w GOARCH=${ARCH} 
	go env -w CGO_ENABLED=1 
	go env -w CC=${CC} 
	go build \
		-buildmode=c-shared \
		-trimpath \
		-ldflags="-s -w -buildid= -H windowsgui" \
		-o ${APP}_${ARCH}.dll 

clean:
	@rm -f ${APP}.* *.dll

.PHONY: ${ARCHES} clean
