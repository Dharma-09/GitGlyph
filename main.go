package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"net/smtp"
)

type GraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		RepoOwner string `json:"repoOwner"`
		RepoName  string `json:"repoName"`
		Label     string `json:"label"`
	} `json:"variables"`
}

type Issue struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Created string `json:"createdAt"`
	Author  struct {
		Login string `json:"login"`
	} `json:"author"`
}

type GraphQLResponse struct {
	Data struct {
		Repository struct {
			Issues struct {
				Edges []struct {
					Node Issue `json:"node"`
				} `json:"edges"`
			} `json:"issues"`
		} `json:"repository"`
	} `json:"data"`
}

type Config struct {
	GithubToken   string
	RepoOwner     string
	RepoName      string
	Label         string
	NotificationEmail string
}

func loadConfig() Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables from the system...")
	}

	config := Config{
		GithubToken:      os.Getenv("GIT_TOKEN"),
		RepoOwner:        os.Getenv("TARGET_REPO_OWNER"),
		RepoName:         os.Getenv("TARGET_REPO_NAME"),
		Label:            os.Getenv("LABEL_TO_TRACK"),
		NotificationEmail: os.Getenv("NOTIFICATION_EMAIL"),
	}

	// Validate required fields
	if config.GithubToken == "" || config.RepoOwner == "" || config.RepoName == "" || config.Label == "" {
		log.Fatalf("Missing required configuration. Please check your .env file or environment variables.")
	}

	return config
}

func fetchIssues(repoOwner, repoName, label string, token string) ([]Issue, error) {
	query := `
		query ($repoOwner: String!, $repoName: String!, $label: String!) {
			repository(owner: $repoOwner, name: $repoName) {
				issues(labels: [$label], states: OPEN, first: 10) {
					edges {
						node {
							title
							url
							createdAt
							author {
								login
							}
						}
					}
				}
			}
		}
	`
	
	// Set up the GraphQL request body
	reqBody := GraphQLRequest{
		Query: query,
	}
	reqBody.Variables.RepoOwner = repoOwner
	reqBody.Variables.RepoName = repoName
	reqBody.Variables.Label = label
	
	// Marshal request into JSON
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request with the access token in the Authorization header
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+loadConfig().GithubToken)
	req.Header.Add("Content-Type", "application/json")

	// Make the request to the GitHub GraphQL API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Log the entire response body for debugging purposes
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	fmt.Println("GraphQL Response: ", string(bodyBytes)) // Add this line to debug the response

	// Decode the response JSON
	var response GraphQLResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract issues from the response
	var issues []Issue
	for _, edge := range response.Data.Repository.Issues.Edges {
		issues = append(issues, edge.Node)
	}
	return issues, nil
}

func sendNotification(issue Issue, recipientEmail string) error {
	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := "SENDER_MAIL" // Replace with your email
	password := "PASSWORD_MAIL"      // Replace with your email password or app-specific password

	// Email content
	subject := fmt.Sprintf("New Issue: %s", issue.Title)
	body := fmt.Sprintf(
		"New issue detected:\n\nTitle: %s\nAuthor: %s\nURL: %s\n\n",
		issue.Title, issue.Author.Login, issue.URL,
	)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// Authentication
	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("Notification sent to %s for issue: %s\n", recipientEmail, issue.Title)
	return nil
}

func main() {
		// Load configuration
		config := loadConfig()

		// Fetch issues
		issues, err := fetchIssues(config.RepoOwner, config.RepoName, config.Label, config.GithubToken)
		if err != nil {
			log.Fatalf("Error fetching issues: %v", err)
		}
	
		// Send notifications for each issue
		for _, issue := range issues {
			err := sendNotification(issue, config.NotificationEmail)
			if err != nil {
				log.Printf("Failed to send notification for issue %s: %v", issue.Title, err)
			}
		}
}
