BASE=dummy
OUT=${BASE}.dll
BUILD=go build
LDFLAGS=-ldflags="-s -w -buildid="
MODE=-buildmode=c-shared
# will enable once go1.13 drops
# TRIM=

.PHONY: windows32 windows64

all: windows32 windows64

clean:
	@rm -f ${BASE}.* *.dll

windows32:
	GOOS=windows \
	GOARCH=386 \
	CGO_ENABLED=1 \
	CC=i686-w64-mingw32-gcc \
	${BUILD} ${MODE} -o 32_${OUT} ${LDFLAGS}

windows64:
	GOOS=windows \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc \
	${BUILD} ${MODE} -o 64_${OUT} ${LDFLAGS}


