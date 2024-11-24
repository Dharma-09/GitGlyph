package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitglyph/functions/github"
)

// WebhookPayload represents the GitHub webhook payload
type WebhookPayload struct {
	Action string `json:"action"`
	Issue  struct {
		Title  string `json:"title"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
		URL string `json:"html_url"`
	} `json:"issue"`
}

// Handle handles incoming GitHub webhook requests
func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var payload WebhookPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Example: Log the webhook event
	log.Printf("Received action: %s, Issue: %s\n", payload.Action, payload.Issue.Title)

	// Handle "labeled" action and check for specific labels
	if payload.Action == "labeled" {
		for _, label := range payload.Issue.Labels {
			if label.Name == "good first issue" {
				log.Printf("Good First Issue found: %s\n", payload.Issue.URL)
				github.PostComment(payload.Issue.URL, "This is a good first issue!")
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"message": "Webhook received"}`)
}
