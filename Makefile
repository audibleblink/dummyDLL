OUT=dummy.dll
BUILD=go build
LDFLAGS=-ldflags="-s -w -buildid="
MODE=-buildmode=c-shared
# will enable once go1.13 drops
# TRIM=

.PHONY: linux windows
linux:
	GOOS=windows \
	CGO_ENABLED=1 \
	CC=i686-w64-mingw32-gcc \
	GOARCH=386 \
	${BUILD} ${MODE} -o ${OUT} ${LDFLAGS}

windows:
	${BUILD} ${MODE} ${TRIM} ${LDFLAGS} -o ${OUT}
