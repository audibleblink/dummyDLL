OUT=dummy.dll
BUILD=go build
LDFLAGS=-ldflags="-s -w -buildid="
MODE=-buildmode=c-shared
# will enable once go1.13 drops
# TRIM=

.PHONY: linux64 linux32 windows64 windows32 windows linux

linux: linux32 linux64
windows: windows32 windows64

linux32:
	GOOS=windows \
	GOARCH=386 \
	CGO_ENABLED=1 \
	CC=i686-w64-mingw32-gcc \
	${BUILD} ${MODE} -o 32_${OUT} ${LDFLAGS}

linux64:
	GOOS=windows \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc \
	${BUILD} ${MODE} -o 64_${OUT} ${LDFLAGS}

windows32:
	GOOS=386 ${BUILD} ${MODE} ${TRIM} ${LDFLAGS} -o ${OUT}

windows32:
	GOOS=amd64 ${BUILD} ${MODE} ${TRIM} ${LDFLAGS} -o ${OUT}
