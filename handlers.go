package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)



func CreateUserHandler(w http.ResponseWriter, r *http.Request){
	userData := User{}
	json.NewDecoder(r.Body).Decode(&userData)
	err := DB.QueryRow(r.Context(), "INSERT INTO users (name, role) VALUES ($1, $2) RETURNING id", userData.Name, userData.Role).Scan(&userData.ID)
	if err != nil {
		http.Error(w, "Failed to insert user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userData)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request){
	id := mux.Vars(r)["id"]
	userData := User{}
	err := DB.QueryRow(r.Context(),
		"SELECT id, name, role FROM users WHERE id=$1",
		id,
	).Scan(&userData.ID, &userData.Name, &userData.Role)

	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userData)	
}