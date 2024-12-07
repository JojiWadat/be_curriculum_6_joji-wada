package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePosts(t *testing.T) {
	// Test GET method
	req, _ := http.NewRequest("GET", "/posts", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePosts)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /posts failed: expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Test POST method
	newPost := Post{
		Content: "Test Post",
		Likes:   0,
		Replies: []Reply{},
		LikedBy: []string{},
		Author: Author{
			UID:  "123",
			Name: "Test User",
		},
	}
	body, _ := json.Marshal(newPost)
	req, _ = http.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /posts failed: expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	// Verify the response body
	var postResponse Post
	err := json.NewDecoder(rr.Body).Decode(&postResponse)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
	if postResponse.Content != newPost.Content {
		t.Errorf("POST /posts content mismatch: expected %s, got %s", newPost.Content, postResponse.Content)
	}
}
