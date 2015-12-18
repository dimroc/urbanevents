package server

import (
	"github.com/dimroc/urbanevents/cityservice/citylib"
	"github.com/labstack/echo"
	"log"
	"time"
	//. "github.com/dimroc/urbanevents/cityservice/utils"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct{}

// Bind attaches api routes
func (api *API) Bind(group *echo.Group) {
	group.Get("/v1/conf", api.ConfHandler)
	group.Get("/v1/settings", api.SettingsHandler)
	group.Get("/v1/cities", api.CitiesHandler)
	group.Get("/v1/cities/:city", api.CityHandler)
	group.Get("/v1/cities/:city/search", api.CitySearchHandler)
}

// ConfHandler handle the app config, for example
func (api *API) ConfHandler(c *echo.Context) error {
	app := c.Get("app").(*App)
	<-time.After(time.Millisecond * 500)
	c.JSON(200, app.Conf.Root)
	return nil
}

func (api *API) SettingsHandler(c *echo.Context) error {
	c.JSON(200, getSettings(c))
	return nil
}

func (api *API) CitiesHandler(c *echo.Context) error {
	settings := getSettings(c)
	elastic := getElasticConnection(c)
	c.JSON(200, settings.GetCities(elastic))
	return nil
}

func (api *API) CityHandler(c *echo.Context) error {
	cityKey := c.Param("city")
	settings := getSettings(c)
	city := settings.FindCity(cityKey)
	elastic := getElasticConnection(c)
	c.JSON(200, city.GetDetails(elastic))
	return nil
}

func (api *API) CitySearchHandler(c *echo.Context) error {
	cityKey := c.Param("city")
	query := c.Query("q")

	settings := getSettings(c)
	elastic := getElasticConnection(c)

	city := settings.FindCity(cityKey)
	c.JSON(200, city.Query(elastic, query))
	return nil
}

func getSettings(c *echo.Context) citylib.Settings {
	rv := c.Get(citylib.CTX_SETTINGS_KEY)
	if rv == nil {
		log.Panic("Could not retrieve Settings")
		return citylib.Settings{}
	} else {
		return rv.(citylib.Settings)
	}
}

func getElasticConnection(c *echo.Context) *citylib.ElasticConnection {
	rv := c.Get(citylib.CTX_ELASTIC_CONNECTION_KEY).(*citylib.ElasticConnection)
	if rv == nil {
		log.Panic("Could not retrieve ElasticConnection")
		return nil
	} else {
		return rv
	}
}
