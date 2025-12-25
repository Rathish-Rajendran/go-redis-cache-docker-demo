package main

import (
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/health", func(w http.ResponseWriter,r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	http.ListenAndServe(port, nil)
}