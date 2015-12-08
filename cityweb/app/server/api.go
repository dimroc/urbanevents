package server

import (
  "log"
  "github.com/gin-gonic/gin"
  "github.com/dimroc/urbanevents/cityservice/citylib"
  //. "github.com/dimroc/urbanevents/cityservice/utils"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct{}

// Bind attaches api routes
func (api *API) Bind(group *gin.RouterGroup) {
	group.GET("/v1/conf", api.ConfHandler)
	group.GET("/v1/settings", api.SettingsHandler)
	group.GET("/v1/cities", api.CitiesHandler)
}

// ConfHandler handle the app config, for example
func (api *API) ConfHandler(c *gin.Context) {
	app := c.MustGet("app").(*App)
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
