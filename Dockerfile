############# builder
FROM golang:1.27.4@sha256:5f1b98fa02fe396ebce2a0227a0e08309e11792915b9dd2e1f71053d7e7e6666 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-os-ubuntu
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make install

############# gardener-extension-os-ubuntu
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS gardener-extension-os-ubuntu
WORKDIR /

COPY --from=builder /go/bin/gardener-extension-os-ubuntu /gardener-extension-os-ubuntu
ENTRYPOINT ["/gardener-extension-os-ubuntu"]
