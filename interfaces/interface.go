package interfaces

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`    //Optional (Handle login in API)
	Avatar    string `json:"avatar"` //Optional (Handle logic in API)
}

type CreatePostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	Pictures    string `json:"pictures"`
	City        string `json:"city"`
	UserID      string `json:"userId"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostID  string `json:"postId"`
	UserID  string `json:"userId"`
}

type LikeRequest struct {
	PostID string `json:"postId"`
	UserID string `json:"userId"`
}

type FollowRequest struct {
	FollowerID  string `json:"followerId"`
	FollowingID string `json:"followingId"`
}

type TagRequest struct {
	PostID string `json:"postId"`
	UserID string `json:"userId"`
}

type CreateResturantRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
}
