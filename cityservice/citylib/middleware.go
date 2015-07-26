package citylib

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"net/http"
)

func ElasticMiddleware(e Elastic) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, CTX_ELASTIC_CONNECTION_KEY, e)
		next(w, r)
	})
}

func SettingsMiddleware(settings Settings) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, CTX_SETTINGS_KEY, settings)
		next(w, r)
	})
}
