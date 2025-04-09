package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type User struct {
	GitHubUsername string `json:"github_username"`
	WrikeEmail     string `json:"wrike_email"`
	WrikeToken     string `json:"wrike_token"`
}

func main() {
	// Get the Base64-encoded credentials from the environment variable
	encodedCredentials := os.Getenv("WRIKE_ME_UP")
	if encodedCredentials == "" {
		log.Fatal("missing WRIKE_ME_UP environment variable")
	}

	// Decode the Base64 string
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		log.Fatalf("failed to decode WRIKE_ME_UP: %v", err)
	}

	// Parse the decoded JSON
	var users []User
	err = json.Unmarshal(decodedCredentials, &users)
	if err != nil {
		log.Fatalf("failed to parse WRIKE_ME_UP: %v", err)
	}

	// Get the GitHub username from the environment variable
	githubUsername := os.Getenv("GITHUB_USERNAME")
	if githubUsername == "" {
		log.Fatal("missing GITHUB_USERNAME environment variable")
	}

	// Find the user with the matching GitHub username
	var user *User
	for _, u := range users {
		if u.GitHubUsername == githubUsername {
			user = &u
		}
	}
	if user == nil {
		log.Fatalf("no credentials found for GitHub user: %s", githubUsername)
	}

	// Get the task ID from the environment variable
	taskID := os.Getenv("WRIKE_TASK_ID")
	if taskID == "" {
		log.Fatal("Missing WRIKE_TASK_ID environment variable")
	}

	// Debug call
	log.Printf("Task ID: %s, Email: %s", taskID, user.WrikeEmail)

	// Create an HTTP client
	client := &http.Client{}

	// Create a new request with the Bearer token
	req, err := http.NewRequest("GET", fmt.Sprintf("https://app-eu.wrike.com/api/v4/tasks/%s/timelogs", taskID), nil)
	if err != nil {
		log.Fatal("failed to create request:", err)
	}
	req.Header.Set("Authorization", "Bearer "+user.WrikeToken)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("wrike API call failed:", err)
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Fatal("error closing response body:", err)
		}
	}(resp.Body)

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body) // Read the response body for error details
		log.Fatalf("Wrike API returned an error: %s (status code: %d)", string(body), resp.StatusCode)
	}

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error reading response body:", err)
	}

	fmt.Println(string(body))
}
