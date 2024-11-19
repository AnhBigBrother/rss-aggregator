package main

import (
	"fmt"
	"net/http"

	"github.com/AnhBigBrother/rss-aggregator/internal/auth"
	"github.com/AnhBigBrother/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			responseErr(w, 401, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			responseErr(w, 400, fmt.Sprintf("couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}