# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/dimroc/urbanevents/cityservice

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

WORKDIR /go/src/github.com/dimroc/urbanevents/cityservice

RUN wget https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm && \
      chmod +x gpm && \
      mv gpm /usr/local/bin # GPM

RUN gpm
RUN go install github.com/dimroc/urbanevents/cityservice

ENTRYPOINT /go/bin/cityservice

EXPOSE 58080
