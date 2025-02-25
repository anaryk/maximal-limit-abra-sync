FROM golang:1.23-alpine as build-stage

RUN apk --no-cache add \
    g++ \
    git \
    make \
    bash

ARG VERSION
ENV VERSION=${VERSION}

WORKDIR /src
COPY . .
RUN go build -o /src/bin/maximal-limit-abra-sync

FROM alpine:latest
RUN apk --no-cache add \
    ca-certificates
COPY --from=build-stage /src/bin/maximal-limit-abra-sync /usr/local/bin/maximal-limit-abra-sync
ENTRYPOINT ["/usr/local/bin/maximal-limit-abra-sync"]