package cityrecorder

import (
	. "github.com/dimroc/urban-events/cityservice/utils"
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"net"
	"net/url"
	"os"
)

var (
	IndexName = os.Getenv("GO_ENV") + "-ntc-geoevents"
)

type ElasticConnection struct {
	Connection *elastigo.Conn
}

func NewElasticConnection(elasticsearchUrl string) *ElasticConnection {
	if len(os.Getenv("GO_ENV")) == 0 {
		log.Panic("GO_ENV not set")
	}

	if len(elasticsearchUrl) == 0 {
		log.Panic("elasticsearchUrl empty")
	} else {
		Logger.Debug("Using Elasticsearch URL " + elasticsearchUrl)
	}

	u, err := url.Parse(elasticsearchUrl)
	if err != nil {
		log.Panic(err)
	}

	host, port, _ := net.SplitHostPort(u.Host)

	connection := elastigo.NewConn()
	connection.Domain = host
	connection.Port = port
	log.Println("Connecting to Elasticsearch", connection.Domain, connection.Port, IndexName)

	return &ElasticConnection{Connection: connection}
}

func (e *ElasticConnection) Write(g GeoEvent) error {
	_, err := e.Connection.Index(IndexName, "tweet", g.Id, nil, g)
	return err
}

func (e *ElasticConnection) Search(query string) elastigo.SearchResult {
	out, err := e.Connection.Search(IndexName, "tweet", nil, query)
	if err != nil {
		log.Panic(err)
	}

	return out
}
