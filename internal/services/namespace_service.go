package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

type NamespaceService struct {
	ScriptPath string
}

func NewNamespaceService() *NamespaceService {
	return &NamespaceService{
		ScriptPath: filepath.Join("scripts", "namespace_handler.py"),
	}
}

func (s *NamespaceService) HandleNamespace(chatID, userID string) (map[string]interface{}, error) {
	cmd := exec.Command("python3", s.ScriptPath, chatID, userID)
	output, err := cmd.CombinedOutput()
	
	// Check command execution error FIRST
	if err != nil {
		return nil, fmt.Errorf("script execution failed: %w", err)
	}

	// Process output after verifying command ran successfully
	cleanedOutput := bytes.TrimSpace(output)
	log.Printf("Namespace handler raw output: %q", cleanedOutput)

	var result map[string]interface{}
	if err := json.Unmarshal(cleanedOutput, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %v\nRaw output: %s", err, cleanedOutput)
	}

	if status, ok := result["status"]; ok && status == "error" {
		return result, fmt.Errorf("namespace handler error: %v", result["message"])
	}

	return result, nil
}
