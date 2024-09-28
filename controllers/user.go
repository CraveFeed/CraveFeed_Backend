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
	"time"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userIdRequest interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&userIdRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("Error decoding request body:", err)
		return
	}
	cachedData, err := Caching.FetchCachedData()
	if err != nil {
		http.Error(w, "Cannot fetch cached data", http.StatusInternalServerError)
		fmt.Println("Error fetching cached data:", err)
		return
	}
	var posts []interfaces.PostData
	err = json.Unmarshal(cachedData, &posts)
	if err != nil {
		http.Error(w, "Error decoding cached data", http.StatusInternalServerError)
		fmt.Println("Error decoding cached data:", err)
		return
	}
	userResponse, err := GetUserData(userIdRequest.Id)
	if err != nil {
		http.Error(w, "Cannot fetch user data", http.StatusInternalServerError)
		fmt.Println("Error fetching user data:", err)
		return
	}
	response := struct {
		User  interfaces.UserResponse `json:"user"`
		Posts []interfaces.PostData   `json:"posts"`
	}{
		User:  userResponse,
		Posts: posts,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}

func GetUserData(userId string) (interfaces.UserResponse, error) {
	pClient := database.PClient
	user, err := pClient.Client.User.FindUnique(
		db.User.ID.Equals(userId),
	).Exec(pClient.Context)
	if err != nil {
		return interfaces.UserResponse{}, fmt.Errorf("cannot fetch user data: %v", err)
	}
	userResponse := interfaces.UserResponse{
		UserId:    user.ID,
		Dish:      user.Dish,
		Type:      user.Type,
		Allergies: user.Allergies,
		Latitude:  user.Latitude,
		Longitude: user.Longitude,
	}
	return userResponse, nil
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userId interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&userId)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("Error decoding request body:", err)
		return
	}
	cachedData, err := Caching.FetchCachedUserData()
	fmt.Println(cachedData)
	if err != nil {
		http.Error(w, "Cannot fetch cached user data", http.StatusInternalServerError)
		fmt.Println("Error fetching cached user data:", err)
		return
	}
	var cachedUsers []interfaces.CachedUser
	err = json.Unmarshal(cachedData, &cachedUsers)
	if err != nil {
		http.Error(w, "Error unmarshalling cached user data", http.StatusInternalServerError)
		fmt.Println("Error unmarshalling cached user data:", err)
		return
	}
	var currentUser interfaces.CachedUser
	var otherUsers []interfaces.CachedUser
	for _, user := range cachedUsers {
		if user.ID == userId.Id {
			currentUser = user
		} else {
			otherUsers = append(otherUsers, user)
		}
	}
	type Response struct {
		CurrentUser interfaces.CachedUser   `json:"currentUser"`
		OtherUsers  []interfaces.CachedUser `json:"otherUsers"`
	}

	response := Response{
		CurrentUser: currentUser,
		OtherUsers:  otherUsers,
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling response data", http.StatusInternalServerError)
		fmt.Println("Error marshalling response data:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
		db.User.Spiciness.Set(userData.Spiciness),
		db.User.Sweetness.Set(userData.Sweetness),
		db.User.Sourness.Set(userData.Sourness),
		db.User.Dish.Set(userData.Dish),
		db.User.Type.Set(userData.Type),
		db.User.Allergies.Set(userData.Allergies),
		db.User.City.Set(userData.City),
		db.User.Longitude.Set(userData.Longitude),
		db.User.Latitude.Set(userData.Latitude),
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
		db.Post.Cuisine.Set(postData.Cuisine),
		db.Post.Dish.Set(postData.Dish),
		db.Post.Type.Set(postData.Type),
		db.Post.Spiciness.Set(postData.Spiciness),
		db.Post.Sweetness.Set(postData.Sweetness),
		db.Post.Sourness.Set(postData.Sourness),
		db.Post.Pictures.Set(postData.Pictures),
		db.Post.UserID.Set(postData.UserID),
		db.Post.City.Set(postData.City),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, createdPost)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient

	var requestBody struct {
		UserID string `json:"userId"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.UserID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		fmt.Println("Error decoding user ID:", err)
		return
	}

	posts, err := pClient.Client.Post.FindMany(
		db.Post.Not(
			db.Post.UserID.Equals(requestBody.UserID),
		),
	).With(
		db.Post.Comments.Fetch().With(
			db.Comment.User.Fetch(),
		).Take(3),
		db.Post.Likes.Fetch(),
		db.Post.RepostedPosts.Fetch(),
		db.Post.User.Fetch(),
	).Exec(pClient.Context)

	if err != nil {
		http.Error(w, "Cannot fetch posts", http.StatusInternalServerError)
		fmt.Println("Error fetching posts:", err)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "No posts found", http.StatusNotFound)
		return
	}

	var responsePosts []interfaces.Post1
	for _, post := range posts {
		var comments []interfaces.Comment
		for _, comment := range post.Comments() {
			commentTimeSinceCreated := time.Since(comment.CreatedAt)
			var commentTimeDescription string
			if commentTimeSinceCreated.Hours() < 24 {
				commentTimeDescription = "a while ago"
			} else {
				commentTimeDescription = fmt.Sprintf("%d days ago", int(commentTimeSinceCreated.Hours()/24))
			}
			comments = append(comments, interfaces.Comment{
				CommentID:   comment.ID,
				Content:     comment.Content,
				UserID:      comment.UserID,
				UserAvatar:  comment.User().Avatar,
				Name:        comment.User().FirstName + " " + comment.User().LastName,
				CommentTime: commentTimeDescription,
			})
		}

		timeSinceUpdated := time.Since(post.UpdatedAt)
		var timeDescription string
		if timeSinceUpdated.Hours() < 24 {
			timeDescription = "a while ago"
		} else {
			timeDescription = fmt.Sprintf("%d days ago", int(timeSinceUpdated.Hours()/24))
		}

		likesCount := len(post.Likes())
		repostsCount := len(post.RepostedPosts())

		responsePosts = append(responsePosts, interfaces.Post1{
			PostID:          post.ID,
			Title:           post.Title,
			Description:     post.Description,
			Longitude:       post.Longitude,
			Latitude:        post.Latitude,
			Pictures:        post.Pictures,
			City:            post.City,
			UserID:          post.UserID,
			UserAvatar:      post.User().Avatar,
			Username:        post.User().Username,
			Name:            post.User().FirstName + " " + post.User().LastName,
			Cuisine:         post.Cuisine,
			Dish:            post.Dish,
			Type:            post.Type,
			Spiciness:       post.Spiciness,
			Sweetness:       post.Sweetness,
			Sourness:        post.Sourness,
			Comments:        comments,
			Likes:           likesCount,
			Reposts:         repostsCount,
			TimeDescription: timeDescription,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.Marshal(responsePosts)
	if err != nil {
		http.Error(w, "Error marshalling response data", http.StatusInternalServerError)
		fmt.Println("Error marshalling response data:", err)
		return
	}
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}

func CheckUserCredentials(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	pClient := database.PClient
	user, err := pClient.Client.User.FindFirst(
		db.User.Username.Equals(credentials.Username),
		db.User.Password.Equals(credentials.Password),
	).Exec(pClient.Context)

	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		fmt.Println("Error fetching user:", err)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	response := struct {
		Message string `json:"message"`
		UserID  string `json:"userId"`
	}{
		Message: "User authenticated successfully",
		UserID:  user.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling response data", http.StatusInternalServerError)
		fmt.Println("Error marshalling response data:", err)
		return
	}
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}

func GetPostsByUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var userId interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&userId)
	if err != nil || userId.Id == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		fmt.Println("Error decoding user ID:", err)
		return
	}
	posts, err := pClient.Client.Post.FindMany(
		db.Post.User.Where(db.User.ID.Equals(userId.Id)),
	).With(
		db.Post.Comments.Fetch().With(
			db.Comment.User.Fetch(),
		).Take(3),
		db.Post.Likes.Fetch(),
		db.Post.RepostedPosts.Fetch(),
		db.Post.User.Fetch(),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "Cannot fetch posts", http.StatusInternalServerError)
		fmt.Println("Error fetching posts:", err)
		return
	}
	if len(posts) == 0 {
		http.Error(w, "No posts found for the user", http.StatusNotFound)
		return
	}
	var responsePosts []interfaces.Post1
	for _, post := range posts {
		var comments []interfaces.Comment
		for _, comment := range post.Comments() {
			commentTimeSinceCreated := time.Since(comment.CreatedAt)
			var commentTimeDescription string
			if commentTimeSinceCreated.Hours() < 24 {
				commentTimeDescription = "a while ago"
			} else {
				commentTimeDescription = fmt.Sprintf("%d days ago", int(commentTimeSinceCreated.Hours()/24))
			}
			comments = append(comments, interfaces.Comment{
				CommentID:   comment.ID,
				Content:     comment.Content,
				UserID:      comment.UserID,
				UserAvatar:  comment.User().Avatar,
				Name:        comment.User().FirstName + " " + comment.User().LastName,
				CommentTime: commentTimeDescription,
			})
		}

		timeSinceUpdated := time.Since(post.UpdatedAt)
		var timeDescription string
		if timeSinceUpdated.Hours() < 24 {
			timeDescription = "a while ago"
		} else {
			timeDescription = fmt.Sprintf("%d days ago", int(timeSinceUpdated.Hours()/24))
		}
		likesCount := len(post.Likes())
		repostsCount := len(post.RepostedPosts())
		responsePosts = append(responsePosts, interfaces.Post1{
			PostID:          post.ID,
			Title:           post.Title,
			Description:     post.Description,
			Longitude:       post.Longitude,
			Latitude:        post.Latitude,
			Pictures:        post.Pictures,
			City:            post.City,
			UserID:          post.UserID,
			UserAvatar:      post.User().Avatar, // Include user's avatar
			Username:        post.User().Username,
			Name:            post.User().FirstName + " " + post.User().LastName,
			Cuisine:         post.Cuisine,
			Dish:            post.Dish,
			Type:            post.Type,
			Spiciness:       post.Spiciness,
			Sweetness:       post.Sweetness,
			Sourness:        post.Sourness,
			Comments:        comments,
			Likes:           likesCount,
			Reposts:         repostsCount,
			TimeDescription: timeDescription,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.Marshal(responsePosts)
	if err != nil {
		http.Error(w, "Error marshalling response data", http.StatusInternalServerError)
		fmt.Println("Error marshalling response data:", err)
		return
	}

	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var userId interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&userId)
	if err != nil || userId.Id == "" {
		http.Error(w, "Invalid request body or missing user ID", http.StatusBadRequest)
		return
	}
	followers, err := pClient.Client.Follows.FindMany(
		db.Follows.FollowingID.Equals(userId.Id),
	).With(
		db.Follows.Follower.Fetch(),
	).Exec(pClient.Context)

	if err != nil {
		http.Error(w, "Cannot fetch followers", http.StatusInternalServerError)
		fmt.Println("Error fetching followers:", err)
		return
	}
	var followerUsers []map[string]string
	for _, follow := range followers {
		followerUsers = append(followerUsers, map[string]string{
			"FirstName": follow.Follower().FirstName,
			"Username":  follow.Follower().Username,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.Marshal(followerUsers)
	if err != nil {
		http.Error(w, "Error marshalling response data", http.StatusInternalServerError)
		fmt.Println("Error marshalling response data:", err)
		return
	}
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pClient := database.PClient
	var userId interfaces.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&userId)
	if err != nil || userId.Id == "" {
		http.Error(w, "Invalid request body or missing user ID", http.StatusBadRequest)
		return
	}
	following, err := pClient.Client.Follows.FindMany(
		db.Follows.FollowerID.Equals(userId.Id),
	).With(
		db.Follows.Following.Fetch(),
	).Exec(pClient.Context)
	if err != nil {
		http.Error(w, "Cannot fetch following users", http.StatusInternalServerError)
		fmt.Println("Error fetching following users:", err)
		return
	}
	var followingUsers []map[string]string
	for _, follow := range following {
		followingUsers = append(followingUsers, map[string]string{
			"FirstName": follow.Following().FirstName,
			"Username":  follow.Following().Username,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.Marshal(followingUsers)
	if err != nil {
		http.Error(w, "Error marshalling response data", http.StatusInternalServerError)
		fmt.Println("Error marshalling response data:", err)
		return
	}
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		fmt.Println("Error writing response:", err)
		return
	}
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
		"noOfPosts":     len(profile.Posts()),
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
		"spiciness":     profile.Spiciness,
		"sweetness":     profile.Sweetness,
		"sourness":      profile.Sourness,
		"type":          profile.Type,
		"allergies":     profile.Allergies,
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
		db.Post.Type.Set(originalPost.Type),
		db.Post.Cuisine.Set(originalPost.Cuisine),
		db.Post.Dish.Set(originalPost.Dish),
		db.Post.Spiciness.Set(originalPost.Spiciness),
		db.Post.Sweetness.Set(originalPost.Sweetness),
		db.Post.Sourness.Set(originalPost.Sourness),
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
	var postData interfaces.EditPostRequest
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	if postData.PostID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	_, err = pClient.Client.Post.FindUnique(
		db.Post.ID.Equals(postData.PostID),
	).Exec(pClient.Context)
	if err != nil {
		log.Printf("Error finding post: %v", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	updatedPost, err := pClient.Client.Post.FindUnique(
		db.Post.ID.Equals(postData.PostID),
	).Update(
		db.Post.Title.Set(postData.Title),
		db.Post.Description.Set(postData.Description),
		db.Post.Longitude.Set(postData.Longitude),
		db.Post.Latitude.Set(postData.Latitude),
		db.Post.Pictures.Set(postData.Pictures),
		db.Post.City.Set(postData.City),
		db.Post.Cuisine.Set(postData.Cuisine),
		db.Post.Dish.Set(postData.Dish),
		db.Post.Type.Set(postData.Type),
		db.Post.Spiciness.Set(postData.Spiciness),
		db.Post.Sweetness.Set(postData.Sweetness),
		db.Post.Sourness.Set(postData.Sourness),
	).Exec(pClient.Context)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	helpers.WriteJSON(w, http.StatusOK, updatedPost)
}
