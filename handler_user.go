package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnhBigBrother/rss-aggregator/internal/auth"
	"github.com/AnhBigBrother/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:name`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		errResponse(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		errResponse(w, 400, fmt.Sprintf("Could not create user: %v", err))
	}
	jsonResponse(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		errResponse(w, 401, fmt.Sprintf("Auth error: %v", err))
	}
	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		errResponse(w, 400, fmt.Sprintf("couldn't get user: %v", err))
		return
	}
	jsonResponse(w, 200, databaseUserToUser(user))
}
