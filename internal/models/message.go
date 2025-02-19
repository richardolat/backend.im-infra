package models

type GitMessage struct {
    UserID     string `json:"userId"`
    ChatID     string `json:"chatId"`
    RepoURL    string `json:"repoURL"`
    CommitHash string `json:"commitHash"`
}

type Response struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}
