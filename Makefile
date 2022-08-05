B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

build: info
	- cd app && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.revision=$(REV) -s -w" -o ../dist/detect_bot

race_test:
	cd app && go test -v -race -mod=vendor -timeout=120s -count 1 ./...

info:
	- @echo "revision $(REV)"