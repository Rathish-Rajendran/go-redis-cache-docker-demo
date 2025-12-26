package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){
	InitRedisClient()
	InitDbClient()

	ctx := context.Background()
	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter,r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	r.HandleFunc("/test-redis", func(w http.ResponseWriter, r *http.Request){
		err := redisClient.Set(ctx,"hello-message","Greetings from Redis!!!",0).Err()
		if err != nil{
			w.Write([]byte("Cannot insert into redis"))
			return
		}

		val, err := redisClient.Get(ctx, "hello-message").Result()
		if err != nil {
			w.Write([]byte("cannot fetch from redis: " + err.Error()))
			return
		}
		w.Write([]byte("Value from redis: " + val))
	}).Methods("GET")

	r.HandleFunc("/users", CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", GetUserHandler).Methods("GET")

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	http.ListenAndServe(port, r)
}