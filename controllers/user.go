package controllers

import (
	database "cravefeed_backend/database"
	helpers "cravefeed_backend/helper"
	"cravefeed_backend/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	pClient := database.PClient
	allUsers, err := pClient.Client.User.FindMany().Exec(pClient.Context)
	if err != nil {
		fmt.Println("Cannot fetch users")
		return

	}
	usersMap := make(map[string]interface{})
	usersMap["users"] = allUsers
	err = helpers.WriteJSON(w, http.StatusOK, usersMap)
	if err != nil {
		fmt.Println("Cannot form response")
		return
	}

}

func GetName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var nameReq interfaces.NameRequest
	err := json.NewDecoder(r.Body).Decode(&nameReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{
		"Name": nameReq.Name,
	})
}
