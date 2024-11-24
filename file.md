## Octobox


## use graphql
query {
  repository(owner: "repository-owner", name: "repository-name") {
    issues(labels: "good first issue", states: OPEN, first: 10) {
      edges {
        node {
          title
          url
        }
      }
    }
  }
}

## outh2
## Web based prompt
## 
curl -X POST https://6db9-173-212-113-180.ngrok-free.app/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "action": "labeled",
    "issue": {
      "title": "Test Issue",
      "labels": [{"name": "good first issue"}],
      "html_url": "https://github.com/Dharma-09/CI-CD-pipelines-examples/issues",
      "number": 123
    },
    "repository": {
      "name": "CI-CD-pipelines-examples",
      "owner": {
        "login": "dharma-09"
      }
    }
  }