package server

import (
	"html/template"
	"io"
	"net/http"

	"github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/itsjamie/go-bindata-templates"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nu7hatch/gouuid"
	"github.com/olebedev/config"
	"os"
)

var (
	settingsFilename = GetenvOrDefault("CITYSERVICE_SETTINGS", "config/cityweb.json")
)

// App struct.
// There is no singleton anti-pattern,
// all variables defined locally inside
// this struct.
type App struct {
	Engine *echo.Echo
	Conf   *config.Config
	React  *React
	API    *API
}

// NewApp returns initialized struct
// of main server application.
func NewApp(opts ...AppOptions) *App {
	options := AppOptions{}
	for _, i := range opts {
		options = i
		break
	}
	options.init()

	// Parse config yaml string from ./conf.go
	conf, err := config.ParseYaml(confString)
	Must(err)
	// Choise a config section by given string
	conf, err = conf.Get(options.Config)
	Must(err)

	// Parse environ variables for defined
	// in config constants
	conf.Env()

	// Make an engine
	engine := echo.New()

	// Set up echo
	engine.SetDebug(conf.UBool("debug"))

	// Regular middlewares
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())

	// Initialize the application
	app := &App{
		Conf:   conf,
		Engine: engine,
		API:    &API{},
		React: NewReact(
			conf.UString("duktape.path"),
			conf.UBool("debug"),
			engine,
		),
	}

	// Use precompiled embedded templates
	app.Engine.SetRenderer(NewTemplate())

	// Map app struct to access from request handlers
	// and middlewares
	app.Engine.Use(func(c *echo.Context) error {
		c.Set("app", app)
		return nil
	})

	// Map uuid for every requests
	app.Engine.Use(func(c *echo.Context) error {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
		return nil
	// Assign settings
	settings, settingsErr := citylib.LoadSettings(settingsFilename)
	Check(settingsErr)
	app.Engine.Use(func(c *gin.Context) {
		c.Set(citylib.CTX_SETTINGS_KEY, settings)
	})

	// Assign Elasticsearch Connection
	elastic := citylib.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	elastic.SetRequestTracer(RequestTracer)
	app.Engine.Use(func(c *gin.Context) {
		c.Set(citylib.CTX_ELASTIC_CONNECTION_KEY, elastic)
	})

	// Avoid favicon react handling
	app.Engine.Get("/favicon.ico", func(c *echo.Context) error {
		c.Redirect(301, "/static/images/favicon.ico")
		return nil
	})

	// Bind api hadling for URL api.prefix
	app.API.Bind(
		app.Engine.Group(
			app.Conf.UString("api.prefix"),
		),
	)

	// Create file http server from bindata
	fileServerHandler := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	})

	// Serve static via bindata and handle via react app
	// in case when static file was not found
	app.Engine.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// execute echo handlers chain
			err := h(c)
			// if page(handler) for url/method not found
			if err != nil && err.Error() == http.StatusText(http.StatusNotFound) {
				// check if file exists
				// omit first `/`
				if _, err := Asset(c.Request().URL.Path[1:]); err == nil {
					fileServerHandler.ServeHTTP(c.Response(), c.Request())
					return nil
				}
				// if static file not found handle request via react application
				return app.React.Handle(c)
			}
			// Move further if err is not `Not Found`
			return err
		}
	})

	return app
}

// Run runs the app
func (app *App) Run() {
	app.Engine.Run(":" + app.Conf.UString("port"))
}

// Custom renderer for Echo, to render html from bindata
type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	return &Template{
		templates: binhtml.New(Asset, AssetDir).MustLoadDirectory("templates"),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// AppOptions is options struct
type AppOptions struct {
	Config string
}

func (ao *AppOptions) init() {
	if ao.Config == "" {
		ao.Config = "local"
	}
}
