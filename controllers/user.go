package controllers

import (
	"cravefeed_backend/Redis/Caching"
	database "cravefeed_backend/database"
	"cravefeed_backend/helper"
	"cravefeed_backend/interfaces"
	"cravefeed_backend/prisma/db"
	"encoding/json"
	"fmt"
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
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	profile, err := pClient.Client.User.FindUnique(
		db.User.ID.Equals(profileData.Id),
	).With(
		db.User.Posts.Fetch(),
		db.User.Followers.Fetch(),
		db.User.Following.Fetch(),
	).Exec(pClient.Context)

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
	defer r.Body.Close()
	pClient := database.PClient
	var profileData interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		"userPosts":     profile.Posts,
		"followers":     profile.Following(),
		"following":     profile.Followers(),
	}

	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, response)
}
