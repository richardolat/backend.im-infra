package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

type TestService struct {
	ScriptPath string
}

func NewTestService() *TestService {
	return &TestService{
		ScriptPath: filepath.Join("scripts", "test-runner.py"),
	}
}

func (s *TestService) RunTests(namespace, repoURL, commit string) (map[string]interface{}, error) {
	cmd := exec.Command("python3", s.ScriptPath,
		"-n", namespace,
		"-r", repoURL,
		"-c", commit,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Test execution failed. Raw output:\n%s", string(output))
		return nil, fmt.Errorf("test runner error: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse test results: %w", err)
	}

	return result, nil
}
