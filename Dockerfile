############# builder
FROM golang:1.14.7 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-os-ubuntu
COPY . .
RUN make install-requirements && make generate && make install

############# gardener-extension-os-ubuntu
FROM alpine:3.12.0 AS gardener-extension-os-ubuntu

COPY --from=builder /go/bin/gardener-extension-os-ubuntu /gardener-extension-os-ubuntu
ENTRYPOINT ["/gardener-extension-os-ubuntu"]
