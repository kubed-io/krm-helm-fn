ARG BUILDER_IMAGE=golang:1.24-alpine
ARG BASE_IMAGE=alpine:3.19

FROM --platform=$BUILDPLATFORM $BUILDER_IMAGE AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /usr/local/bin/function ./cmd/function

# Install Helm
RUN apk add --no-cache curl
ARG HELM_VERSION="v3.13.2"
RUN curl -fsSL -o /helm-${HELM_VERSION}-${TARGETOS}-${TARGETARCH}.tar.gz https://get.helm.sh/helm-${HELM_VERSION}-${TARGETOS}-${TARGETARCH}.tar.gz && \
     tar -zxvf /helm-${HELM_VERSION}-${TARGETOS}-${TARGETARCH}.tar.gz && \
     mv ./${TARGETOS}-${TARGETARCH}/helm /usr/local/bin/helm

# Final image
FROM $BASE_IMAGE AS final
RUN apk add --no-cache ca-certificates git
COPY --from=builder /usr/local/bin/function /usr/local/bin/function
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm

ENV PATH=/usr/local/bin:$PATH \
    HELM_PATH_CACHE=/var/cache \
    HELM_CONFIG_HOME=/tmp/helm/config \
    HELM_CACHE_HOME=/tmp/helm/cache

ENTRYPOINT ["function"]

FROM mcr.microsoft.com/devcontainers/go:1.24 as devcontainer

# install kustomize
RUN <<EOF
set -e
ARCH="${TARGETPLATFORM#*/}"
# Handle common architecture names
case $ARCH in
  "amd64") KUSTOMIZE_ARCH="amd64" ;;
  "arm64") KUSTOMIZE_ARCH="arm64" ;;
  "arm/v7") KUSTOMIZE_ARCH="arm" ;;
  *) KUSTOMIZE_ARCH="amd64" ;;
esac

KUSTOMIZE_VERSION="v5.7.1"
echo "installing Kustomize ${KUSTOMIZE_VERSION} for ${KUSTOMIZE_ARCH}"

# Download and verify the URL works, then extract
curl -fsSL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F${KUSTOMIZE_VERSION}/kustomize_${KUSTOMIZE_VERSION}_linux_${KUSTOMIZE_ARCH}.tar.gz" | tar -xz -C /usr/local/bin
chmod +x /usr/local/bin/kustomize

# Verify installation
kustomize version || echo "Warning: kustomize installation may have issues"
EOF
