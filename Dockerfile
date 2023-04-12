FROM golang:1.16-alpine as backend

ARG GIT_BRANCH
ARG GITHUB_SHA
ARG CI

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN apk add --no-cache --update git tzdata ca-certificates

RUN \
    version=${GIT_BRANCH}-${GITHUB_SHA:0:7}-$(date +%Y%m%dT%H:%M:%S) && \
    echo "version=$version" && \
    cd app && go build -o /build/cam_cp -ldflags "-X main.revision=${version} -s -w"


FROM ghcr.io/umputun/baseimage/app:v1.6.1 as base

FROM scratch

COPY --from=backend /build/cam_cp /app/cam_cp
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

WORKDIR /app
ENTRYPOINT ["/app/cam_cp"]