package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnhBigBrother/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		Name string `json:name`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Could not create feed: %v", err))
		return
	}
	responseJson(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}
	responseJson(w, 201, databaseFeedsToFeeds(feeds))
}
