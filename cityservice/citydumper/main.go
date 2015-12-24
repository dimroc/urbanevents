package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	elastigo "github.com/dimroc/elastigo/lib"
	citylib "github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "citydumper"
	app.Version = "0.0.1"
	app.Usage = "Dump a city's geoevents from elasticsearch to JSONL."

	var citykey, after, before, filename, elasticsearchUrl string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "citykey, c",
			Usage:       "Required. The key for the city you are trying to dump (aka export)",
			Destination: &citykey,
			EnvVar:      "CITYKEY",
		},
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "http://localhost:9200",
			Usage:       "The name of the file to write to",
			EnvVar:      "ELASTICSEARCH_URL",
			Destination: &elasticsearchUrl,
		},
		cli.StringFlag{
			Name:        "after, a",
			Value:       "1980-01-01",
			Usage:       "The date string which geoevents must be after",
			Destination: &after,
		},
		cli.StringFlag{
			Name:        "before, b",
			Value:       time.Now().Format("2006-01-02 03:04:05 Z0700"),
			Usage:       "The date string which geoevents must be before",
			Destination: &before,
		},
		cli.StringFlag{
			Name:        "output, o",
			Value:       "/tmp/citydump.jsonl",
			Usage:       "The name of the file to write to",
			Destination: &filename,
		},
	}

	app.Action = func(c *cli.Context) {
		if len(citykey) == 0 {
			fmt.Println("Must set citykey flag. See help")
			return
		}

		fmt.Println("Dumping city " + citykey + " from " + elasticsearchUrl)
		elastic := citylib.NewElasticConnection(elasticsearchUrl)
		elastic.SetRequestTracer(RequestTracer)
		outputfile, err := os.Create(filename)
		Check(err)
		defer outputfile.Close()
		defer elastic.Close()

		// Set up Output File
		writer := bufio.NewWriter(outputfile)
		defer writer.Flush()

		// Set up Elasticsearch reading
		dsl := elastigo.Search(citylib.ES_IndexName).Type(citylib.ES_TypeName).Size("1000").
			Pretty().Filter(elastigo.Filter().And(
			elastigo.Filter().Term("city", citykey),
			elastigo.Filter().Range("createdAt", after, nil, before, nil, ""),
		),
		)

		searchResult := elastic.ScanAndScrollDsl(*dsl)
		Logger.Debug("Scroll ID: " + searchResult.ScrollId)

		for {
			searchResult = elastic.Scroll(searchResult.ScrollId)

			Logger.Debug("Scroll ID: " + searchResult.ScrollId)
			Logger.Debug(searchResult.String())

			geoevents := citylib.GeoEventsFromElasticSearch(&searchResult)
			for _, geoevent := range geoevents {
				writeGeoevent(writer, geoevent)
			}

			if len(geoevents) == 0 {
				break
			}
		}

		fmt.Println("Output written to " + outputfile.Name())
	}

	app.Run(os.Args)
}

func writeGeoevent(writer *bufio.Writer, geoevent citylib.GeoEvent) {
	_, err := writer.WriteString(ToJsonStringUnsafe(geoevent) + "\n")
	Check(err)
}
