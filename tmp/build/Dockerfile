FROM alpine:3.6

RUN adduser -D kubernetes-pod-restarter
USER kubernetes-pod-restarter

ADD tmp/_output/bin/kubernetes-pod-restarter /usr/local/bin/kubernetes-pod-restarter
