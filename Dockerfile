############# builder
FROM golang:1.23.6 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-os-ubuntu
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make install

############# gardener-extension-os-ubuntu
FROM gcr.io/distroless/static-debian11:nonroot AS gardener-extension-os-ubuntu
WORKDIR /

COPY --from=builder /go/bin/gardener-extension-os-ubuntu /gardener-extension-os-ubuntu
ENTRYPOINT ["/gardener-extension-os-ubuntu"]
