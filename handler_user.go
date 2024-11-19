package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		responseErr(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Could not create user: %v", err))
		return
	}
	responseJson(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseJson(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetPostFollowed(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostFollowed(r.Context(), database.GetPostFollowedParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}
	responseJson(w, 200, databasePostsToPosts(posts))
}
