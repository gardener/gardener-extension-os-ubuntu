############# builder
FROM eu.gcr.io/gardener-project/3rd/golang:1.15.5 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-os-ubuntu
COPY . .
RUN make install-requirements && make generate && make install

############# gardener-extension-os-ubuntu
FROM eu.gcr.io/gardener-project/3rd/alpine:3.12.1 AS gardener-extension-os-ubuntu

COPY --from=builder /go/bin/gardener-extension-os-ubuntu /gardener-extension-os-ubuntu
ENTRYPOINT ["/gardener-extension-os-ubuntu"]
