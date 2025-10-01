ARG BUILDER_IMAGE=golang:1.21-alpine
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
FROM $BASE_IMAGE
RUN apk add --no-cache ca-certificates git
COPY --from=builder /usr/local/bin/function /usr/local/bin/function
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm

ENV PATH /usr/local/bin:$PATH
ENV HELM_PATH_CACHE /var/cache
ENV HELM_CONFIG_HOME /tmp/helm/config
ENV HELM_CACHE_HOME /tmp/helm/cache

ENTRYPOINT ["function"]