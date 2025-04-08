# GitGlyph
[![Netlify Status](https://api.netlify.com/api/v1/badges/9ab4ddfc-9591-456e-aff5-75167167d0cf/deploy-status)](https://app.netlify.com/sites/gitglyph/deploys)
[![GitHub stars](https://img.shields.io/github/stars/Dharma-09/GitGlyph)](https://github.com/Dharma-09/GitGlyph/stargazers)
![Last Commit](https://img.shields.io/github/last-commit/Dharma-09/GitGlyph)
![Pull Requests](https://img.shields.io/github/issues-pr-raw/Dharma-09/GitGlyph)
[![License](https://img.shields.io/github/license/Dharma-09/GitGlyph)](https://github.com/Dharma-09/GitGlyph/blob/master/LICENSE)

# GitGlyph

GitGlyph is a tool to track issues in any GitHub repository with a specific label. It sends notifications whenever a new issue with the specified label is detected...

---

## Features
- Track issues in any GitHub repository by specifying the repository and label.
- Notifications are sent to a specified email address (e.g., via SMTP).
- Easy setup with environment variables.

---

## Getting Started

### Prerequisites
- Go (1.20 or later)
- A GitHub Personal Access Token with `repo` and `read:org` permissions.
- An SMTP server (e.g., Gmail) for notifications (optional).

---

### Setup Instructions

1. **Fork this Repository**

   - Click the "Fork" button at the top-right of this page to create your own copy of this repository.
OR
1. **Clone the Repository**

   ```git
   git clone https://github.com/<your-username>/GitGlyph.git
   cd GitGlyph
   ```
2. **Configure Environment Variables**

    Copy the .env.example file to .env:
    ```bash
    cp .env.example .env
    ```
    Edit the .env file with your preferred values:
    ```plaintext
    GITHUB_TOKEN=your_personal_access_token
    TARGET_REPO_OWNER=target_repo_owner
    TARGET_REPO_NAME=target_repo_name
    LABEL_TO_TRACK=label_to_track
    NOTIFICATION_EMAIL=your_email@example.com
    ```
3. Run the Application

    ```go
    go run main.go
    ```
