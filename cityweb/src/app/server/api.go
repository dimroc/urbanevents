package server

import (
	"github.com/dimroc/urbanevents/cityservice/citylib"
	"log"
	//. "github.com/dimroc/urbanevents/cityservice/utils"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct{}

// Bind attaches api routes
func (api *API) Bind(group *echo.Group) {
	group.Get("/v1/conf", api.ConfHandler)
	group.GET("/v1/settings", api.SettingsHandler)
	group.GET("/v1/cities", api.CitiesHandler)
	group.GET("/v1/cities/:city", api.CityHandler)
	group.GET("/v1/cities/:city/search", api.CitySearchHandler)
}

// ConfHandler handle the app config, for example
func (api *API) ConfHandler(c *echo.Context) error {
	app := c.Get("app").(*App)
	<-time.After(time.Millisecond * 500)
	c.JSON(200, app.Conf.Root)
}

func (api *API) SettingsHandler(c *gin.Context) {
	c.JSON(200, getSettings(c))
}

func (api *API) CitiesHandler(c *gin.Context) {
	settings := getSettings(c)
	elastic := getElasticConnection(c)
	c.JSON(200, settings.GetCityDetails(elastic))
}

func (api *API) CityHandler(c *gin.Context) {
	cityKey := c.Param("city")
	settings := getSettings(c)
	city := settings.FindCity(cityKey)
	elastic := getElasticConnection(c)
	c.JSON(200, city.GetDetails(elastic))
}

func (api *API) CitySearchHandler(c *gin.Context) {
	cityKey := c.Param("city")
	query := c.DefaultQuery("q", "")

	settings := getSettings(c)
	elastic := getElasticConnection(c)

	city := settings.FindCity(cityKey)
	c.JSON(200, city.Query(elastic, query))
}

func getSettings(c *gin.Context) citylib.Settings {
	rv := c.MustGet(citylib.CTX_SETTINGS_KEY)
	if rv == nil {
		log.Panic("Could not retrieve Settings")
		return citylib.Settings{}
	} else {
		return rv.(citylib.Settings)
	}
}

func getElasticConnection(c *gin.Context) *citylib.ElasticConnection {
	rv := c.MustGet(citylib.CTX_ELASTIC_CONNECTION_KEY).(*citylib.ElasticConnection)
	if rv == nil {
		log.Panic("Could not retrieve ElasticConnection")
		return nil
	} else {
		return rv
	}
}
