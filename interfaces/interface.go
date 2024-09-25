package interfaces

type CreateUserRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`    //Optional (Handle login in API)
	Avatar    string `json:"avatar"` //Optional (Handle logic in API)
	Spiciness int    `json:"spiciness"`
	Sweetness int    `json:"sweetness"`
	Sourness  int    `json:"sourness"`
	Type      string `json:"type"`
	Allergies string `json:"allergies"`
	City      string `json:"city"`
}

type CreateProfileRequest struct {
	Id string `json:"id"`
}

type CreateProfileIdRequest struct {
	PostId string `json:"postId"`
}

type CreateUsernameRequest struct {
	Username string `json:"username"`
}

type RepostRequest struct {
	PostID string `json:"postId"`
	UserID string `json:"userId"`
}

type CreatePostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	Pictures    string `json:"pictures"`
	City        string `json:"city"`
	UserID      string `json:"userId"`
	Cuisine     string `json:"cuisine"`
	Dish        string `json:"dish"`
	Type        string `json:"type"`
	Spiciness   int    `json:"spiciness"`
	Sweetness   int    `json:"sweetness"`
	Sourness    int    `json:"sourness"`
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

type EditPostRequest struct {
	PostID      string `json:"postId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	Pictures    string `json:"pictures"`
	City        string `json:"city"`
	Cuisine     string `json:"cuisine"`
	Dish        string `json:"dish"`
	Type        string `json:"type"`
	Spiciness   int    `json:"spiciness"`
	Sweetness   int    `json:"sweetness"`
	Sourness    int    `json:"sourness"`
}

type Comment struct {
	CommentID string `json:"commentId"`
	Content   string `json:"content"`
	UserID    string `json:"userId"`
}

type Post struct {
	PostID      string    `json:"postId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Longitude   string    `json:"longitude"`
	Latitude    string    `json:"latitude"`
	Pictures    string    `json:"pictures"`
	City        string    `json:"city"`
	UserID      string    `json:"userId"`
	Cuisine     string    `json:"cuisine"`
	Dish        string    `json:"dish"`
	Type        string    `json:"type"`
	Spiciness   int       `json:"spiciness"`
	Sweetness   int       `json:"sweetness"`
	Sourness    int       `json:"sourness"`
	Comments    []Comment `json:"comments"`
	Likes       int       `json:"likes"`
	Reposts     int       `json:"reposts"`
}
