package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	"github.com/gorilla/context"
	"net/http"
)

func ElasticMiddleware(e cityrecorder.Elastic) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, CTX_ELASTIC_CONNECTION_KEY, e)
		next(w, r)
	})
}

func SettingsMiddleware(settings cityrecorder.Settings) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, CTX_SETTINGS_KEY, settings)
		next(w, r)
	})
}
