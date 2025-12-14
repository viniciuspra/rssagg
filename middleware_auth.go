package main

import (
	"fmt"
	"net/http"

	"github.com/viniciuspra/rssagg/internal/auth"
	"github.com/viniciuspra/rssagg/internal/db"
)

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
