package macro

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadMacroFromFile(filePath string) ([]Action, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read macro file: %w", err)
	}

	var actions []Action
	err = json.Unmarshal(fileBytes, &actions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal macro file: %w", err)
	}

	return actions, nil
}