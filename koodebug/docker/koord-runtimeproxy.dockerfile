FROM --platform=$TARGETPLATFORM golang:1.18 as builder
WORKDIR /go/src/github.com/koordinator-sh/koordinator

ARG VERSION
ARG TARGETARCH
ENV VERSION $VERSION
ENV GOOS linux
ENV GOARCH $TARGETARCH

COPY go.mod go.mod
COPY go.sum go.sum
RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w GO111MODULE=on
RUN go mod download

COPY apis/ apis/
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 go build -a -o koord-runtime-proxy cmd/koord-runtime-proxy/main.go

FROM --platform=$TARGETPLATFORM centos:7
WORKDIR /
COPY --from=builder /go/src/github.com/koordinator-sh/koordinator/koord-runtime-proxy .
ENTRYPOINT ["/koord-runtime-proxy"]
