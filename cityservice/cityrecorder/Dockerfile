FROM alpine:3.2
RUN apk add --update ca-certificates
ADD tmp/cityrecorder /go/bin/cityrecorder
ADD . /go/src/github.com/dimroc/urbanevents/cityservice/cityrecorder
ADD tmp/config /go/src/github.com/dimroc/urbanevents/cityservice/cityrecorder

WORKDIR /go/src/github.com/dimroc/urbanevents/cityservice/cityrecorder
ENTRYPOINT /go/bin/cityrecorder
