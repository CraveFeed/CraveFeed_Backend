package main

import (
	"cravefeed_backend/database"
	router "cravefeed_backend/routers"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
}

func (app *Application) Serve() error {
	port := app.Config.Port
	fmt.Printf("Serving app on port %s", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router.Routes(),
	}
	return srv.ListenAndServe()
}

func main() {
	db, err := database.ConnectDB()

	if err != nil {
		fmt.Println("Database cannot be connected")
	}

	defer func() {
		if db.Client != nil {
			db.Client.Disconnect()
		}
	}()

	config := Config{
		Port: "3000",
	}

	app := &Application{
		Config: config,
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}

}
