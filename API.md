# WebSocket API Specification

## Endpoint
`ws://yourdomain.com/ws`

## Message Format
```json
{
    "userId": "string",
    "chatId": "string", 
    "repoURL": "string",
    "commitHash": "string"
}
```

## Response Types
| Type       | Description                          | Example                      |
|------------|--------------------------------------|------------------------------|
| `success`  | Operation completed                  | `{type: "success", payload:}`|
| `error`    | Error response                       | `{type: "error", message:}`  |
| `progress` | Real-time operation updates          | `{type: "progress", stage:}` |

## Error Codes
- 1000: Invalid message format
- 1001: Repository access error
- 1002: Git operation failed
- 1003: Invalid commit hash

## Example Usage
```json
// Request
{
    "userId": "user123",
    "chatId": "chat456",
    "repoURL": "https://github.com/example/repo",
    "commitHash": "a1b2c3d4"
}

// Success Response
{
    "type": "success",
    "payload": {
        "repoPath": "/app/repos/repo",
        "commit": "a1b2c3d4"
    }
}
```
