FROM golang:1.10-alpine as builder

RUN apk add --no-cache \
    bash \
    curl \
    git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/rebuy-de/kubernetes-pod-restarter
COPY . .

RUN dep ensure -vendor-only -v
RUN ./tmp/build/build.sh


FROM alpine:3.6

RUN adduser -D kubernetes-pod-restarter
USER kubernetes-pod-restarter

COPY --from=builder /go/src/github.com/rebuy-de/kubernetes-pod-restarter/tmp/_output/bin/kubernetes-pod-restarter /usr/local/bin/kubernetes-pod-restarter
