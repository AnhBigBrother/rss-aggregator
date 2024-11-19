package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnhBigBrother/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Could not create feed-follow: %v", err))
		return
	}
	responseJson(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollowed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedsFollowed, err := apiCfg.DB.GetFeedFollowed(r.Context(), user.ID)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("couldn't get feed followed: %v", err))
		return
	}
	responseJson(w, 200, databaseFeedsToFeeds(feedsFollowed))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIdStr := chi.URLParam(r, "feedId")
	feedId, err := uuid.Parse(feedIdStr)
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedId,
		UserID: user.ID,
	})
	if err != nil {
		responseErr(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}
	responseJson(w, 200, struct {
		Message string `json:"message"`
	}{Message: "Deleted successfully"})
}
