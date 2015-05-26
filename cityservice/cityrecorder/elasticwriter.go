package cityrecorder

import (
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"net"
	"net/url"
	"os"
)

var (
	IndexName = os.Getenv("GO_ENV") + "-ntc-geoevents"
)

type ElasticWriter struct {
	Connection *elastigo.Conn
}

func NewElasticWriter(elasticsearchUrl string) Writer {
	if len(os.Getenv("GO_ENV")) == 0 {
		log.Fatal("GO_ENV not set")
	}

	u, err := url.Parse(elasticsearchUrl)
	if err != nil {
		log.Fatal(err)
	}

	host, port, _ := net.SplitHostPort(u.Host)

	connection := elastigo.NewConn()
	connection.Domain = host
	connection.Port = port
	log.Println("Connecting to Elasticsearch", connection)

	return &ElasticWriter{Connection: connection}
}

func (e *ElasticWriter) Write(g GeoEvent) error {
	_, err := e.Connection.Index(IndexName, "tweet", g.Id, nil, g)
	return err
}
