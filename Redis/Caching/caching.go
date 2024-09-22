package Caching

import (
	"context"
	"cravefeed_backend/Redis"
	"cravefeed_backend/database"
	"encoding/json"
	"fmt"
	"time"
)

var ctx = context.Background()

const cacheKey = "all_posts"

func fetchAndCachePosts() error {
	pClient := database.PClient
	allPosts, err := pClient.Client.Post.FindMany().Exec(pClient.Context)
	if err != nil {
		return fmt.Errorf("cannot fetch posts: %v", err)
	}
	postsJSON, err := json.Marshal(allPosts)
	if err != nil {
		return fmt.Errorf("cannot serialize posts: %v", err)
	}
	rdb := Redis.GetClient()
	err = rdb.Set(ctx, cacheKey, postsJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("cannot cache posts: %v", err)
	}
	return nil
}

func UpdateCachePeriodically() {
	for {
		err := fetchAndCachePosts()
		if err != nil {
			fmt.Println("Error updating cache:", err)
		} else {
			fmt.Println("Cache updated successfully")
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
	return cachedData, nil
}
