# Docker image for Drone's downstream trigger plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-downstream .

FROM alpine:3.2
RUN apk add -U ca-certificates && rm -rf /var/cache/apk/*
ADD drone-downstream /bin/
ENTRYPOINT ["/bin/drone-downstream"]
