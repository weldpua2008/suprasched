FROM golang:1.15-alpine AS build-env

WORKDIR /go/src/github.com/weldpua2008/suprasched

COPY . .

RUN apk --no-cache add build-base git mercurial gcc alpine-sdk

# docker image build --no-cache --tag suprasched:snapshot --build-arg GIT_COMMIT=$(git log -1 --format=%H) .
ARG GIT_COMMIT

RUN go get -d -v ./... && \
    go install -v ./... && \
    export GIT_COMMIT_LOACAL=$(git log -1 --format=%H) && \
    export GIT_COMMIT=${GIT_COMMIT:-$GIT_COMMIT_LOACAL} && \
    go build -o /root/suprasched -ldflags="-X github.com/weldpua2008/suprasched/cmd.GitCommit=${GIT_COMMIT}" main.go

FROM alpine:latest

# Bring in metadata via --build-arg
ARG BRANCH=unknown
ARG IMAGE_CREATED=unknown
ARG IMAGE_REVISION=unknown
ARG IMAGE_VERSION=unknown

# Configure image labels
LABEL \
    # https://github.com/opencontainers/image-spec/blob/master/annotations.md
    branch=$branch \
    maintainer="Valeriy Soloviov" \
    org.opencontainers.image.authors="Valeriy Soloviov" \
    org.opencontainers.image.created=$IMAGE_CREATED \
    org.opencontainers.image.description="The orchestration layer around jobs, observe the execution time, rescheduler and to control concurrent execution." \
    org.opencontainers.image.documentation="https://github.com/weldpua2008/suprasched/" \
    org.opencontainers.image.licenses="Apache License 2.0" \
    org.opencontainers.image.revision=$IMAGE_REVISION \
    org.opencontainers.image.source="https://github.com/weldpua2008/suprasched/" \
    org.opencontainers.image.title="suprasched" \
    org.opencontainers.image.url="https://github.com/weldpua2008/suprasched/" \
    org.opencontainers.image.vendor="suprasched" \
    org.opencontainers.image.version=$IMAGE_VERSION

# Default image environment variable settings
ENV org.opencontainers.image.created=$IMAGE_CREATED \
    org.opencontainers.image.revision=$IMAGE_REVISION \
    org.opencontainers.image.version=$IMAGE_VERSION


WORKDIR /root/

# Copy source
COPY --from=build-env /root/suprasched .

RUN adduser -D --shell /bin/bash hadoop && \
    apk --no-cache add curl bash

# Set entrypoint
ENTRYPOINT ["/root/suprasched"]
