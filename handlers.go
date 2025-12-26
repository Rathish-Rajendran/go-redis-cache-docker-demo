package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
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
	keyInRedis := "user:"+id

	valFromRedis, err := redisClient.Get(r.Context(), keyInRedis).Result()
	if err == nil {
		json.Unmarshal([]byte(valFromRedis), &userData)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"source": "cache",
			"data": userData,
		})
		return
	}

	err = DB.QueryRow(r.Context(),
		"SELECT id, name, role FROM users WHERE id=$1",
		id,
	).Scan(&userData.ID, &userData.Name, &userData.Role)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	userDataInJsonBytes, err := json.Marshal(userData)
	if err != nil {
		log.Warnf("error json marshalling user data before inserting into redis for user id:%v", id)
	}
	err = redisClient.Set(r.Context(), keyInRedis, userDataInJsonBytes, 30*time.Second).Err()
	if err != nil {
		log.Warnf("error inserting data to redis for user id:%v", id)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"source": "database",
		"data": userData,
	})
}