# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/dimroc/urbanevents/cityservice

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

WORKDIR /go/src/github.com/dimroc/urbanevents/cityservice

RUN go get
RUN go install github.com/dimroc/urbanevents/cityservice

ENTRYPOINT /go/bin/cityservice

EXPOSE 58080
