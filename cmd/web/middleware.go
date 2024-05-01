package main

import "net/http"

func (app *Config) SessionLoad(next http.Handler) http.Handler {
	// The LoadAndSave method takes an http.Handler as an argument and returns an http.Handler.
	// LoadAndSave performs tasks related to loading session data from the incoming request and saving it back after the request is handled.
	return app.Session.LoadAndSave(next)
}

func (app *Config) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "userID") {
			app.Session.Put(r.Context(), "warning", "You must log in to see this page")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
