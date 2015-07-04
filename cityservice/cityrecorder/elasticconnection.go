package cityrecorder

import (
	. "github.com/dimroc/urbanevents/cityservice/utils"
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"net"
	"net/url"
	"os"
)

var (
	IndexName = os.Getenv("GO_ENV") + "-geoevents-write"
)

type Elastic interface {
	Writer
	Close()
	Search(query string) elastigo.SearchResult
}

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
		Logger.Debug("Using Elasticsearch URL " + elasticsearchUrl + " with index " + IndexName)
	}

	u, err := url.Parse(elasticsearchUrl)
	if err != nil {
		log.Panic(err)
	}

	connection := elastigo.NewConn()

	if u.User != nil {
		connection.Username = u.User.Username()
		p, _ := u.User.Password()
		connection.Password = p
	}

	host, port, _ := net.SplitHostPort(u.Host)
	connection.Domain = host
	connection.Port = port

	return &ElasticConnection{Connection: connection}
}

func (e *ElasticConnection) Close() {
	e.Connection.Close()
}

func (e *ElasticConnection) SetRequestTracer(requestTracer func(string, string, string)) {
	e.Connection.RequestTracer = requestTracer
}

func (e *ElasticConnection) Refresh() error {
	// Panics if tracing requests during refresh, so temporarily silence requesting
	oldTracer := e.Connection.RequestTracer
	e.Connection.RequestTracer = nil
	_, err := e.Connection.Refresh(IndexName)
	Check(err)
	e.Connection.RequestTracer = oldTracer
	return err
}

func (e *ElasticConnection) Write(g GeoEvent) error {
	_, err := e.Connection.Index(IndexName, "geoevent", g.Id, nil, g)
	return err
}

func (e *ElasticConnection) Search(query string) elastigo.SearchResult {
	out, err := e.Connection.Search(IndexName, "geoevent", nil, query)
	if err != nil {
		log.Panic(err)
	}

	return out
}

type BulkElasticConnection struct {
	*ElasticConnection
	BulkIndexer *elastigo.BulkIndexer
}

func NewBulkElasticConnection(elasticsearchUrl string) *BulkElasticConnection {
	elastic := NewElasticConnection(elasticsearchUrl)
	bulkIndexer := elastic.Connection.NewBulkIndexerErrors(5, 10)
	bulkIndexer.Start()

	return &BulkElasticConnection{elastic, bulkIndexer}
}

func (e *BulkElasticConnection) Write(g GeoEvent) error {
	return e.BulkIndexer.Index(IndexName, "geoevent", g.Id, "", &g.CreatedAt, g, false)
}
