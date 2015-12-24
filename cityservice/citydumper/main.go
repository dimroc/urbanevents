package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	elastigo "github.com/dimroc/elastigo/lib"
	citylib "github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "citydumper"
	app.Version = "0.0.1"
	app.Usage = "Dump a city's geoevents from elasticsearch to CSV."

	var citykey string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "citykey, c",
			Usage:       "The key for the city you are trying to dump (aka export)",
			Destination: &citykey,
			EnvVar:      "CITYKEY",
		},
	}

	app.Action = func(c *cli.Context) {
		if len(citykey) == 0 {
			fmt.Println("Must set citykey flag. See help")
			return
		}

		elastic := citylib.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
		defer elastic.Close()

		dsl := elastigo.Search(citylib.ES_IndexName).Size("1000").
			SearchType("scan").Scroll("20s"). // Scan and scroll: https://www.elastic.co/guide/en/elasticsearch/guide/current/scan-scroll.html
			Type(citylib.ES_TypeName).Pretty().Filter(
			elastigo.Filter().Term("city", citykey),
		).Sort(
			elastigo.Sort("createdAt").Desc(),
		)

		scrollResult := elastic.SearchDsl(*dsl)
		Logger.Debug("Scroll ID: " + scrollResult.ScrollId)

		scrollArgs := map[string]interface{}{"scroll": "20s"}
		searchResult, err := elastic.Connection.Scroll(scrollArgs, scrollResult.ScrollId)
		Check(err)

		Logger.Debug(searchResult.String())
		geoevents := citylib.GeoEventsFromElasticSearch(&searchResult)
		for _, geoevent := range geoevents {
			Logger.Debug(geoevent.String())
		}
	}

	app.Run(os.Args)
}
