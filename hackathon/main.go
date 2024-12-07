package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var client *firestore.Client

type Post struct {
	Content string `json:"content"`
	Author  struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
	} `json:"author"`
}

func main() {
	// Firestoreのクライアントを初期化
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "your-gcp-project-id", option.WithCredentialsFile("path/to/your/credentials.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	http.HandleFunc("/posts", handlePosts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newPost Post
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&newPost); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Firestoreにデータを保存
		_, _, err := client.Collection("posts").Add(r.Context(), newPost)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create post: %v", err), http.StatusInternalServerError)
			return
		}

		// レスポンスを返す
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPost)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
