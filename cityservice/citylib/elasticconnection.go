package citylib

import (
	"fmt"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"net"
	"net/url"
	"os"
)

type Elastic interface {
	Writer
	Close()
	Search(query string) elastigo.SearchResult
	SearchDsl(query elastigo.SearchDsl) *elastigo.SearchResult
	Percolate(geojson GeoJson) []string
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
		Logger.Notice("Using Elasticsearch URL " + elasticsearchUrl + " with index " + ES_IndexName)
	}

	u, err := url.Parse(elasticsearchUrl)
	Check(err)

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
	_, err := e.Connection.Refresh(ES_IndexName)
	Check(err)
	e.Connection.RequestTracer = oldTracer
	return err
}

func (e *ElasticConnection) Write(g GeoEvent) error {
	_, err := e.Connection.Index(ES_IndexName, ES_TypeName, g.Id, nil, g)
	return err
}

func (e *ElasticConnection) Search(query string) elastigo.SearchResult {
	out, err := e.Connection.Search(ES_IndexName, ES_TypeName, nil, query)
	Check(err)
	return out
}

func (e *ElasticConnection) SearchDsl(query elastigo.SearchDsl) *elastigo.SearchResult {
	out, err := query.Result(e.Connection)
	Check(err)
	return out
}

func (e *ElasticConnection) Percolate(geojson GeoJson) []string {
	//func (c *Conn) Percolate(index string, _type string, name string, args map[string]interface{}, doc string) (Match, error) {

	geojson_str := ToJsonStringUnsafe(&geojson)
	doc := fmt.Sprintf(`{"doc": { "geojson": %s }}`, geojson_str)

	result, err := e.Connection.Percolate(ES_IndexName, ES_TypeName, "", nil, doc)
	Check(err)

	hoods := make([]string, len(result.Matches))
	for index, match := range result.Matches {
		hoods[index] = match.Id
	}

	return hoods
}

// Bulk Elastic
// Shit don't work. Silently fails after around a day of usage. gg BulkIndexer.
// Probably related to: https://github.com/mattbaird/elastigo/commit/0c98885a2b2575c99882263dfc4bf6aae9079a63
//type BulkElasticConnection struct {
//*ElasticConnection
//BulkIndexer *elastigo.BulkIndexer
//}

//func NewBulkElasticConnection(elasticsearchUrl string) *BulkElasticConnection {
//elastic := NewElasticConnection(elasticsearchUrl)
//bulkIndexer := elastic.Connection.NewBulkIndexerErrors(5, 10)
//bulkIndexer.Start()

//return &BulkElasticConnection{elastic, bulkIndexer}
//}

//func (e *BulkElasticConnection) Write(g GeoEvent) error {
//return e.BulkIndexer.Index(ES_IndexName, ES_TypeName, g.Id, "", &g.CreatedAt, g, false)
//}
