FROM golang:1.19-alpine AS builder

RUN apk add --no-cache curl tar
ENV K8S_VERSION=1.21.2
# download kubebuilder tools required by envtest
RUN curl -sSLo envtest-bins.tar.gz "https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-${K8S_VERSION}-$(go env GOOS)-$(go env GOARCH).tar.gz"
RUN tar -vxzf envtest-bins.tar.gz -C /usr/local/

WORKDIR /api

COPY . .

RUN CGO_ENABLED=0 go build -o server

FROM alpine
ENV ETCD_UNSUPPORTED_ARCH=arm64
# Add new non-root user 'mohamed'
RUN adduser -D mohamed
USER mohamed
WORKDIR /home/mohamed

# required by api server to determine config/crds path
ENV GOPATH=/go
COPY --from=builder /go/pkg/mod/github.com/mohamed-abdelrhman/pack-dispatch /go/pkg/mod/github.com/mohamed-abdelrhman/pack-dispatch
COPY --from=builder /api/server /home/mohamed/api/server

EXPOSE 5000
ENTRYPOINT [ "./api/server" ]
