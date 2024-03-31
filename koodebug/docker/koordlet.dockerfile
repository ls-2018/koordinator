FROM --platform=$TARGETPLATFORM golang:1.18 as builder
WORKDIR /go/src/github.com/koordinator-sh/koordinator

ARG VERSION
ARG TARGETARCH
ENV VERSION $VERSION
ENV GOOS linux
ENV GOARCH $TARGETARCH

RUN #apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 871920D1991BC93C

COPY koodebug/debian.list /etc/apt/sources.list

RUN apt update && apt install -y bash build-essential cmake wget
RUN wget https://sourceforge.net/projects/perfmon2/files/libpfm4/libpfm-4.13.0.tar.gz && \
  echo "bcb52090f02bc7bcb5ac066494cd55bbd5084e65  libpfm-4.13.0.tar.gz" | sha1sum -c && \
  tar -xzf libpfm-4.13.0.tar.gz && \
  rm libpfm-4.13.0.tar.gz
RUN export DBG="-g -Wall" && \
  make -e -C libpfm-4.13.0 && \
  make install -C libpfm-4.13.0

COPY go.mod go.mod
COPY go.sum go.sum
RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w GO111MODULE=on
RUN go mod download

COPY apis/ apis/
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go build -a -o koordlet cmd/koordlet/main.go

# The CUDA container images provide an easy-to-use distribution for CUDA supported platforms and architectures.
# NVIDIA provides rich images in https://hub.docker.com/r/nvidia/cuda/tags, literally cover all kinds of CUDA version
# and system architecture. Please replace the following base image according to your Kubernetes/System environment.
# For more details about how those images got built, you might wanna check the original Dockerfile in
# https://gitlab.com/nvidia/container-images/cuda/-/tree/master/dist.

FROM --platform=$TARGETPLATFORM nvidia/cuda:11.6.2-base-ubuntu20.04
WORKDIR /
RUN apt update && apt install -y lvm2 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/src/github.com/koordinator-sh/koordinator/koordlet .
COPY --from=builder /usr/local/lib /usr/lib
ENTRYPOINT ["/koordlet"]