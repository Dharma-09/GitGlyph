# GitGlyph
[![Netlify Status](https://api.netlify.com/api/v1/badges/9ab4ddfc-9591-456e-aff5-75167167d0cf/deploy-status)](https://app.netlify.com/sites/gitglyph/deploys)


GitGlyph is a Go-based application that tracks GitHub issues in real-time based on specified repositories and labels. Users can provide access tokens for seamless tracking of public and private repositories without requiring webhook setups.

## Features  
- **Real-time Issue Tracking**: Dynamically fetch GitHub issues based on repository and label.  
- **Secure Access**: Utilize personal access tokens to authenticate and fetch data securely.  
- **Custom Filters**: Specify target repositories and labels for focused tracking.  
- **GraphQL Integration**: Uses the GitHub GraphQL API for efficient and flexible data queries.  

---

## Prerequisites  
Ensure you have the following installed:  
1. [Go](https://go.dev/dl/) (version 1.20 or higher)  
2. A [GitHub Personal Access Token](https://github.com/settings/tokens) with the `repo` and `read:org` permissions.  
3. A terminal or code editor to run Go programs.  

---

## Installation  

1. **Clone the Repository**:  
   ```bash
   git clone https://github.com/Dharma-09/gitglyph.git
   cd gitglyph

   go run main.go
   ```