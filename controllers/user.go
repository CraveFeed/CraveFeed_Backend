package controllers

import (
	"cravefeed_backend/Redis/Caching"
	database "cravefeed_backend/database"
	"cravefeed_backend/helper"
	"cravefeed_backend/interfaces"
	"cravefeed_backend/prisma/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	cachedData, err := Caching.FetchCachedData()
	if err != nil {
		http.Error(w, "Cannot fetch cached data", http.StatusInternalServerError)
		fmt.Println("Error fetching cached data:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(cachedData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) { //Removed optional field in Avatar/Bio for testing
	defer r.Body.Close()
	pClient := database.PClient
	var userData interfaces.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := pClient.Client.User.CreateOne(
		db.User.Email.Set(userData.Email),
		db.User.Username.Set(userData.Username),
		db.User.Password.Set(userData.Password),
		db.User.Bio.Set(userData.Bio),
		db.User.Avatar.Set(userData.Avatar),
		db.User.FirstName.Set(userData.FirstName),
		db.User.LastName.Set(userData.LastName),
	).Exec(pClient.Context)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, createdUser)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var postData interfaces.CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&postData)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdPost, err := pClient.Client.Post.CreateOne(
		db.Post.Title.Set(postData.Title),
		db.Post.Description.Set(postData.Description),
		db.Post.Longitude.Set(postData.Longitude),
		db.Post.Latitude.Set(postData.Latitude),
		db.Post.Pictures.Set(postData.Pictures),
		db.Post.UserID.Set(postData.UserID),
		db.Post.City.Set(postData.City),
	).Exec(pClient.Context)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, createdPost)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var commentData interfaces.CreateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&commentData)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdComment, err := pClient.Client.Comment.CreateOne(
		db.Comment.Content.Set(commentData.Content),
		db.Comment.PostID.Set(commentData.PostID),
		db.Comment.UserID.Set(commentData.UserID),
	).Exec(pClient.Context)

	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, createdComment)
}

func HandleLikeRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var likeReq interfaces.LikeRequest
	err := json.NewDecoder(r.Body).Decode(&likeReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdLike, err := pClient.Client.Like.CreateOne(
		db.Like.PostID.Set(likeReq.PostID),
		db.Like.UserID.Set(likeReq.UserID),
	).Exec(pClient.Context)

	if err != nil {
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, createdLike)
}

func HandleFollowRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var followReq interfaces.FollowRequest
	err := json.NewDecoder(r.Body).Decode(&followReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	createdFollow, err := pClient.Client.Follows.CreateOne(
		db.Follows.FollowerID.Set(followReq.FollowerID),
		db.Follows.FollowingID.Set(followReq.FollowingID),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "Failed to create follow", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, createdFollow)
}

func GetProfileBio(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var profileData interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileData)
	profile, err := pClient.Client.User.FindUnique(
		db.User.ID.Equals(profileData.Id),
	).With(
		db.User.Posts.Fetch(),
		db.User.Followers.Fetch(),
		db.User.Following.Fetch(),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return

	}
	response := map[string]interface{}{
		"username":      profile.Username,
		"bio":           profile.Bio,
		"avatar":        profile.Avatar,
		"firstname":     profile.FirstName,
		"lastname":      profile.LastName,
		"noOfFollowers": len(profile.Followers()),
		"noOfFollowing": len(profile.Following()),
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, response)
}

func GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	pClient := database.PClient
	var profileData interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileData)
	if err != nil || profileData.Id == "" {
		http.Error(w, "Invalid request body or missing user ID", http.StatusBadRequest)
		return
	}
	userID := profileData.Id
	profile, err := pClient.Client.User.FindUnique(
		db.User.ID.Equals(userID),
	).With(
		db.User.Posts.Fetch(),
		db.User.Followers.Fetch(),
		db.User.Following.Fetch(),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	response := map[string]interface{}{
		"username":      profile.Username,
		"bio":           profile.Bio,
		"avatar":        profile.Avatar,
		"firstname":     profile.FirstName,
		"lastname":      profile.LastName,
		"coverImage":    profile.Avatar,
		"noOfPosts":     len(profile.Posts()),
		"noOfFollowers": len(profile.Followers()),
		"noOfFollowing": len(profile.Following()),
		"userPosts":     profile.Posts(),
		"followers":     profile.Followers(),
		"following":     profile.Following(),
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, response)
}

func Repost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var repostData interfaces.RepostRequest
	err := json.NewDecoder(r.Body).Decode(&repostData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := repostData.UserID
	originalPostID := repostData.PostID
	originalPost, err := pClient.Client.Post.FindUnique(
		db.Post.ID.Equals(originalPostID),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "Original post not found", http.StatusNotFound)
		return
	}
	existingRepost, err := pClient.Client.Post.FindFirst(
		db.Post.UserID.Equals(userID),
		db.Post.OriginalPostID.Equals(originalPostID),
	).Exec(pClient.Context)
	if existingRepost != nil {
		http.Error(w, "Repost already exists for this user", http.StatusConflict)
		return
	}
	newPost, err := pClient.Client.Post.CreateOne(
		db.Post.UserID.Set(userID),
		db.Post.Title.Set("Repost: "+originalPost.Title),
		db.Post.Description.Set(originalPost.Description),
		db.Post.Longitude.Set(originalPost.Longitude),
		db.Post.Latitude.Set(originalPost.Latitude),
		db.Post.Pictures.Set(originalPost.Pictures),
		db.Post.City.Set(originalPost.City),
		db.Post.OriginalPostID.Set(originalPostID),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "Failed to create repost", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, newPost)
}

func GetReposts(w http.ResponseWriter, r *http.Request) {
	pClient := database.PClient
	var profileData interfaces.CreateProfileIdRequest
	err := json.NewDecoder(r.Body).Decode(&profileData)
	reposts, err := pClient.Client.Post.FindMany(
		db.Post.OriginalPostID.Equals(profileData.PostId),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "Failed to fetch reposts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, reposts)
}

func GetUsernameUserId(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var profileUsername interfaces.CreateUsernameRequest
	err := json.NewDecoder(r.Body).Decode(&profileUsername)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	username := profileUsername.Username
	profile, err := pClient.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	response := map[string]interface{}{
		"id":       profile.ID,
		"username": profile.Username,
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, response)
}

func EditPosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var profileData interfaces.EditPostRequest
	err := json.NewDecoder(r.Body).Decode(&profileData)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	if profileData.PostID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	posts, err := pClient.Client.Post.UpsertOne(
		db.Post.ID.Equals(profileData.PostID), // Ensure this matches your ORM's method
	).Update(
		db.Post.Title.Set(profileData.Title),
		db.Post.Description.Set(profileData.Description),
		db.Post.Longitude.Set(profileData.Longitude),
		db.Post.Latitude.Set(profileData.Latitude),
		db.Post.Pictures.Set(profileData.Pictures),
		db.Post.City.Set(profileData.City),
	).Exec(pClient.Context)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		http.Error(w, "Failed to Update Post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, posts)
}
