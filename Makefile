include .env
export
B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)

docker:
	docker build -t i39.in/cam_cp:master --progress=plain .
docker_yolo_gpu:
	docker build -t i39.in/cam_cp_gpu:master -f Darknet.dockerfile --progress=plain .
build: info
	- cd app &&  GOARCH=amd64  go build  -ldflags "-X main.revision=$(REV) -s -w " -o ../dist/cam_cp

race_test:
	cd app && go test -v -race -mod=vendor -timeout=120s -count 1 ./...

test:
	cd app && go test -v  -mod=vendor -timeout=120s -count 1 ./...

info:
	- @echo "revision $(REV)"
