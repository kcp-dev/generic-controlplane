# syntax=docker/dockerfile:1.4

# Copyright 2024 The KCP Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the binary
FROM --platform=${BUILDPLATFORM} docker.io/golang:1.22.2 AS builder
WORKDIR /workspace

# Install dependencies.
RUN apt-get update && apt-get install -y jq && mkdir bin

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
USER 0

# Install kubectl.
RUN wget "https://dl.k8s.io/release/$(go list -m -json k8s.io/kubernetes | jq -r .Version)/bin/linux/$(uname -m | sed 's/aarch.*/arm64/;s/armv8.*/arm64/;s/x86_64/amd64/')/kubectl" -O bin/kubectl && chmod +x bin/kubectl

ENV GOPRIVATE=github.com/faroshq/cluster-proxy
ARG GH_TOKEN

# and so that source changes don't invalidate our downloaded layer
ENV GOPROXY=direct

RUN --mount=type=cache,target=/go/pkg/mod \
    git config --global url."https://${GH_TOKEN}:@github.com/".insteadOf "https://github.com/" && \
    go mod download && \
    git config --global --unset url."https://${GH_TOKEN}:@github.com/".insteadOf

# Copy the sources
COPY ./ ./

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    make OS=${TARGETOS} ARCH=${TARGETARCH}

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:debug
WORKDIR /
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder workspace/bin/gcp /
ENV KUBECONFIG=/etc/gcp/config/admin.kubeconfig
# Use uid of nonroot user (65532) because kubernetes expects numeric user when applying pod security policies
RUN ["/busybox/sh", "-c", "mkdir -p /data && chown 65532:65532 /data"]
USER 65532:65532
WORKDIR /data
VOLUME /data
ENTRYPOINT ["/gcp"]
CMD ["start"]
