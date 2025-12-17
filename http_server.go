package main

import (
	"log"
	"net/http"
)

func runHTTPServer(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	  log.Println("http server error:", err)
	}
}
