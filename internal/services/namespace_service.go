package services

import (
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
		ScriptPath: filepath.Join("scripts", "namespace_handler.sh"),
	}
}

func (s *NamespaceService) HandleNamespace(chatID, userID string) (map[string]interface{}, error) {
	// Construct the identifier argument for the script
	identifier := fmt.Sprintf("%s-%s", chatID, userID)
	
	// Make script executable
	if err := exec.Command("chmod", "+x", s.ScriptPath).Run(); err != nil {
		return nil, fmt.Errorf("failed to make script executable: %v", err)
	}

	// Execute the script
	cmd := exec.Command(s.ScriptPath, identifier)
	output, err := cmd.CombinedOutput()
	log.Printf("Namespace handler output: %s", string(output))

	if err != nil {
		return nil, fmt.Errorf("namespace handler failed: %v", err)
	}

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse script output: %v", err)
	}

	return result, nil
}
