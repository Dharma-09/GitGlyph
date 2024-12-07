package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// PostComment posts a comment to a GitHub issue
func PostComment(issueURL, comment string) {
	token := os.Getenv("GIT_TOKEN")
	if token == "" {
		log.Println("GIT_TOKEN not set in environment variables.")
		return
	}

	url := fmt.Sprintf("%s/comments", issueURL)
	payload := map[string]string{"body": comment}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to post comment. Status: %d\n", resp.StatusCode)
	}
}
