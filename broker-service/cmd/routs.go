package main

import "net/http"

func routes() http.Handler {
	mux := chi.NewRouter(

		// specify who is allowed to connect
		mux.
	)
}