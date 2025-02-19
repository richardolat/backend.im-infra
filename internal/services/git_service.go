package services

import (
    "fmt"
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
    cmd := exec.Command(s.ScriptPath, repoURL, commitHash)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("git handler failed: %v, output: %s", err, string(output))
    }
    return nil
}
