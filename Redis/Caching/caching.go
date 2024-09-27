package Caching

import (
	"context"
	"cravefeed_backend/Redis"
	"cravefeed_backend/database"
	"cravefeed_backend/prisma/db"
	"encoding/json"
	"fmt"
	"time"
)

var ctx = context.Background()

const (
	cacheKey  = "all_posts"
	cacheKey1 = "all_users"
)

func fetchAndCachePosts() error {
	pClient := database.PClient
	allPosts, err := pClient.Client.Post.FindMany().With(
		db.Post.Likes.Fetch(),
	).Exec(pClient.Context)
	if err != nil {
		return fmt.Errorf("cannot fetch posts: %v", err)
	}
	type CachedPost struct {
		PostID     string `json:"post_id"`
		Cuisine    string `json:"cuisine"`
		Dish       string `json:"dish"`
		Type       string `json:"type"`
		LikesCount int    `json:"likes_count"`
		Latitude   string `json:"latitude"`
		Longitude  string `json:"longitude"`
		Sweetness  int    `json:"sweetness"`
		Sourness   int    `json:"sourness"`
		Spiciness  int    `json:"spiciness"`
	}
	var cachedPosts []CachedPost
	for _, post := range allPosts {
		cachedPost := CachedPost{
			PostID:     post.ID,
			Cuisine:    post.Cuisine,
			Dish:       post.Dish,
			Type:       post.Type,
			LikesCount: len(post.Likes()),
			Latitude:   post.Latitude,
			Longitude:  post.Longitude,
			Sweetness:  post.Sweetness,
			Sourness:   post.Sourness,
			Spiciness:  post.Spiciness,
		}
		cachedPosts = append(cachedPosts, cachedPost)
	}
	postsJSON, err := json.Marshal(cachedPosts)
	if err != nil {
		return fmt.Errorf("cannot serialize posts: %v", err)
	}
	rdb := Redis.GetClient()
	err = rdb.Set(pClient.Context, cacheKey, postsJSON, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("cannot cache posts: %v", err)
	}
	return nil
}

func fetchAndCacheUsers() error {
	pClient := database.PClient
	allUsers, err := pClient.Client.User.FindMany().With(
		db.User.Followers.Fetch(),
	).Exec(pClient.Context)
	if err != nil {
		return fmt.Errorf("cannot fetch users: %v", err)
	}
	type CachedUser struct {
		ID            string `json:"id"`
		Type          string `json:"type"`
		Dish          string `json:"dish"`
		City          string `json:"city"`
		FollowerCount int    `json:"followerCount"`
	}
	var cachedUsers []CachedUser
	for _, user := range allUsers {
		cachedUser := CachedUser{
			ID:            user.ID,
			Type:          user.Type,
			Dish:          user.Dish,
			City:          user.City,
			FollowerCount: len(user.Followers()),
		}
		cachedUsers = append(cachedUsers, cachedUser)
	}
	usersJSON, err := json.Marshal(cachedUsers)
	if err != nil {
		return fmt.Errorf("cannot serialize users: %v", err)
	}
	rdb := Redis.GetClient()
	err = rdb.Set(pClient.Context, cacheKey1, usersJSON, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("cannot cache users: %v", err)
	}
	return nil
}

func UpdateCachePeriodically() {
	for {
		if err := fetchAndCachePosts(); err != nil {
			fmt.Println("Error updating posts cache:", err)
		} else {
			fmt.Println("Posts cache updated successfully")
		}

		if err := fetchAndCacheUsers(); err != nil {
			fmt.Println("Error updating users cache:", err)
		} else {
			fmt.Println("Users cache updated successfully")
		}

		time.Sleep(10 * time.Second)
	}
}

func FetchCachedData() ([]byte, error) {
	rdb := Redis.GetClient()
	cachedData, err := rdb.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch cached data: %v", err)
	}

	if len(cachedData) == 0 {
		return nil, fmt.Errorf("cached data is empty")
	}

	return cachedData, nil
}

func FetchCachedUserData() ([]byte, error) {
	rdb := Redis.GetClient()
	cachedData, err := rdb.Get(ctx, cacheKey1).Bytes()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch cached user data: %v", err)
	}
	if len(cachedData) == 0 {
		return nil, fmt.Errorf("cached user data is empty")
	}
	return cachedData, nil
}
