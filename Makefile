APP=dummy
FLAGS=-trimpath -ldflags="-s -w -buildid="
MODE=-buildmode=c-shared
ARCH = $(word 1, $@)

all: 386 amd64

clean:
	@rm -f ${APP}.* *.dll

386:
	GOOS=windows \
	GOARCH=${ARCH} \
	CGO_ENABLED=1 \
	CC=i686-w64-mingw32-gcc \
	go build ${MODE} -o ${APP}_${ARCH}.dll ${FLAGS}

amd64:
	GOOS=windows \
	GOARCH=${ARCH} \
	CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc \
	go build ${MODE} -o ${APP}_${ARCH}.dll ${FLAGS}

.PHONY: all 386 amd64 clean
