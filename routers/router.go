package routers

import (
	"cravefeed_backend/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/createUser", controllers.CreateUser)
	router.Post("/createPost", controllers.CreatePost)
	router.Post("/createComment", controllers.CreateComment)
	router.Post("/like", controllers.HandleLikeRequest)
	router.Post("/follow", controllers.HandleFollowRequest)
	router.Post("/posts", controllers.GetAllPosts)
	router.Post("/getProfileInfo", controllers.GetProfileInfo)
	router.Post("/getProfileBio", controllers.GetProfileBio)
	router.Post("/repostPost", controllers.Repost)
	router.Post("/getReposts", controllers.GetReposts)
	router.Post("/searchByUsername", controllers.GetUsernameUserId)
	router.Post("/editPost", controllers.EditPosts)
	router.Post("/getPosts", controllers.GetPosts)
	router.Post("/getPostsById", controllers.GetPostsByUsers)
	router.Post("/getFollowers", controllers.GetFollowers)
	router.Post("/getFollowing", controllers.GetFollowing)
	router.Post("/users", controllers.GetAllUsers)

	return router

}
