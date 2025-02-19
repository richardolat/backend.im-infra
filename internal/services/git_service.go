package services

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
)

type GitService struct {
    ScriptPath string
}

func NewGitService() *GitService {
    return &GitService{
        ScriptPath: filepath.Join("scripts", "git_handler.sh"),
    }
}

func (s *GitService) HandleRepository(repoURL, commitHash string) error {
    // Log the script path and parameters
    log.Printf("Executing git handler script: %s with repo: %s and commit: %s", s.ScriptPath, repoURL, commitHash)

    // Make script executable
    if err := exec.Command("chmod", "+x", s.ScriptPath).Run(); err != nil {
        return fmt.Errorf("failed to make script executable: %v", err)
    }

    // Get absolute path of script
    absPath, err := filepath.Abs(s.ScriptPath)
    if err != nil {
        return fmt.Errorf("failed to get absolute path: %v", err)
    }

    // Execute the script
    cmd := exec.Command(absPath, repoURL, commitHash)
    
    // Get the current working directory
    pwd, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("failed to get working directory: %v", err)
    }
    log.Printf("Current working directory: %s", pwd)

    // Capture both stdout and stderr
    output, err := cmd.CombinedOutput()
    log.Printf("Git handler output: %s", string(output))
    
    if err != nil {
        return fmt.Errorf("git handler failed: %v, output: %s", err, string(output))
    }

    return nil
}
