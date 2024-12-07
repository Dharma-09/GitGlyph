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

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
}

func fetchIssues(repoOwner, repoName, label, token string) ([]Issue, error) {
	query := `
		query($repoOwner: String!, $repoName: String!, $label: String!) {
			repository(owner: $repoOwner, name: $repoName) {
				issues(first: 10, labels: [$label], states: OPEN) {
					nodes {
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
	`

	requestBody := GraphQLRequest{
		Query: query,
	}
	requestBody.Variables.RepoOwner = repoOwner
	requestBody.Variables.RepoName = repoName
	requestBody.Variables.Label = label

	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // Use io.ReadAll instead of ioutil.ReadAll
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GraphQL Response: %s", string(body))
	}

	var response struct {
		Data struct {
			Repository struct {
				Issues struct {
					Nodes []Issue `json:"nodes"`
				} `json:"issues"`
			} `json:"repository"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %v", err)
	}

	return response.Data.Repository.Issues.Nodes, nil
}

func main() {
	// Load environment variables
	loadEnv()

	repoOwner := os.Getenv("REPO_OWNER")
	repoName := os.Getenv("REPO_NAME")
	label := os.Getenv("LABEL")
	token := os.Getenv("GIT_TOKEN")

	if repoOwner == "" || repoName == "" || label == "" || token == "" {
		log.Fatal("Required environment variables are missing")
	}

	// Fetch issues
	issues, err := fetchIssues(repoOwner, repoName, label, token)
	if err != nil {
		log.Fatalf("Error fetching issues: %v", err)
	}

	// Display issues
	for _, issue := range issues {
		fmt.Printf("New issue: %s by %s\n", issue.Title, issue.Author.Login)
		fmt.Printf("URL: %s\n", issue.URL)
	}
}
