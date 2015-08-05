FROM gliderlabs/alpine:3.1

ADD tmp/cityrecorder /go/bin/cityrecorder
ADD . /go/src/github.com/dimroc/urbanevents/cityservice/cityrecorder
ADD tmp/config /go/src/github.com/dimroc/urbanevents/cityservice/cityrecorder

WORKDIR /go/src/github.com/dimroc/urbanevents/cityservice/cityrecorder
ENTRYPOINT /go/bin/cityrecorder
