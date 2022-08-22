############# builder
FROM golang:1.17.13 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-os-ubuntu
COPY . .
RUN make generate && make install

############# gardener-extension-os-ubuntu
FROM gcr.io/distroless/static-debian11:nonroot AS gardener-extension-os-ubuntu
WORKDIR /

COPY --from=builder /go/bin/gardener-extension-os-ubuntu /gardener-extension-os-ubuntu
ENTRYPOINT ["/gardener-extension-os-ubuntu"]
