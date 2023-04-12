ARG CUDA_VERSION=11.7.0
ARG BASE_IMAGE=nvidia/cuda:$CUDA_VERSION-cudnn8-devel-ubuntu20.04

FROM $BASE_IMAGE AS builder

ENV DEBIAN_FRONTEND noninteractive

ARG GOLANG_VERSION=1.19
ENV GOLANG_VERSION=$GOLANG_VERSION


RUN apt-get update \
      && apt-get install --no-install-recommends --no-install-suggests -y gnupg2 ca-certificates \
            git build-essential tzdata curl \            
      && rm -rf /var/lib/apt/lists/*

ARG LATEST_COMMIT=9d40b619756be9521bc2ccd81808f502daaa3e9a
ENV LATEST_COMMIT=$LATEST_COMMIT

# Golang env
ARG GOLANG_VERSION=1.19.2
ENV GOLANG_VERSION=$GOLANG_VERSION
ENV GO_PARENT_DIR /usr
ENV PATH /usr/go/bin:$PATH

RUN git clone https://github.com/AlexeyAB/darknet.git && cd darknet \
      && git checkout $LATEST_COMMIT \
      && sed -i -e 's/GPU=0/GPU=1/g' Makefile \
	  && sed -i -e 's/CUDNN=0/CUDNN=1/g' Makefile \
	  && sed -i -e 's/LIBSO=0/LIBSO=1/g' Makefile \
	  && make -j $(shell nproc --all) \
	  && cp libdarknet.so /usr/lib/libdarknet.so \
      && cp include/darknet.h /usr/include/darknet.h \
      && ldconfig \
      && cd .. && rm -rf darknet \
      && curl -sSf https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz | tar -xz -C "$GO_PARENT_DIR"


ADD . /build
WORKDIR /build

ARG GIT_BRANCH
ARG GITHUB_SHA


ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=1

ENV CUDA_MM_VERSION=11.7

RUN version=$(date +%Y%m%dT%H:%M:%S) \
    && echo "version=$version" \
    && echo $CUDA_MM_VERSION \
    && cd app && go build -o /build/cam_cp -ldflags "-X main.revision=${version} -s -w  -v -r /usr/local/cuda-$CUDA_MM_VERSION/targets/x86_64-linux/lib/:/usr/local/cuda-$CUDA_MM_VERSION/compat/"

 
ARG CUDA_VERSION=11.7.0
FROM nvidia/cuda:$CUDA_VERSION-cudnn8-runtime-ubuntu20.04

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update \
      && apt-get install --no-install-recommends --no-install-suggests -y gnupg2 ca-certificates vlc-bin vlc-plugin-base netcat \            
      && rm -rf /var/lib/apt/lists/*

ENV CUDA_MM_VERSION=11.7
COPY --from=builder /usr/local/cuda-$CUDA_MM_VERSION/compat/libcuda.so.1 /usr/local/cuda/targets/x86_64-linux/lib/libcuda.so.1
COPY --from=builder /usr/lib/libdarknet.so /usr/lib/libdarknet.so
COPY --from=builder /usr/include/darknet.h /usr/include/darknet.h
COPY --from=builder /build/cam_cp /app/cam_cp
COPY ./run.sh /app/run.sh
RUN useradd -m vlcuser

RUN ldconfig

WORKDIR /app
ENTRYPOINT ["/app/run.sh"]
EXPOSE 8080