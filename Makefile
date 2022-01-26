VERSION?=$(shell git describe --abbrev=0)+$(shell date +'%Y%m%d%H%I%S')

EXECUTABLE:=civ

DESTDIR?=
PREFIX?=/usr/local

# BINARIES
# ==============================================================================
all: ${EXECUTABLE}

${EXECUTABLE}:
	GOPROXY=$(shell go env GOPROXY) \
	GOPRIVATE=$(shell go env GOPRIVATE) \
		go build -ldflags "-X main.version=${VERSION:v%=%}" -o ${@}

# UN/INSTALL
# ==============================================================================
PHONY+=install
install: ${EXECUTABLE}
	install --directory ${DESTDIR}${PREFIX}/bin
	install --mode 755 ${EXECUTABLE} ${DESTDIR}${PREFIX}/bin/${EXECUTABLE}

	install --directory ${DESTDIR}${PREFIX}/licenses/${EXECUTABLE}
	install --mode 644 LICENSE ${DESTDIR}${PREFIX}/licenses/${EXECUTABLE}/LICENSE

PHONY+=uninstall
uninstall:
	-rm --recursive --force \
		${DESTDIR}${PREFIX}/bin/${EXECUTABLE} \
		${DESTDIR}${PREFIX}/licenses/${EXECUTABLE}/LICENSE

# CLEAN
# ==============================================================================
PHONY+=clean
clean:
	rm --force --recursive ${EXECUTABLE}* || true

# TEST
# ==============================================================================
PHONY+=test/unit
test/unit:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic -timeout 600s -count=1 ./pkg/...

PHONY+=test/coverage
test/coverage: test/unit
	go tool cover -html=coverage.txt

# GOLANGCI-LINT
# ==============================================================================
PHONY+=golangci-lint
golangci-lint:
	golangci-lint run --concurrency=$(shell nproc)

# GOSEC
# ==============================================================================
PHONY+=gosec
gosec:
	gosec $(shell pwd)/...

# PHONY
# ==============================================================================
# Declare the contents of the PHONY variable as phony.  We keep that information
# in a variable so we can use it in if_changed.
.PHONY: ${PHONY}