package main

import "net/http"

func (app *Config) SessionLoad(next http.Handler) http.Handler {
	// The LoadAndSave method takes an http.Handler as an argument and returns an http.Handler.
	// LoadAndSave performs tasks related to loading session data from the incoming request and saving it back after the request is handled.
	return app.Session.LoadAndSave(next)
}
