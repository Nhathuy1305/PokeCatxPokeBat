package utils

import (
	"encoding/json"
	"os"
)

func ReadJSON(filename string, data interface{}) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileData, data)
}

func WriteJSON(filename string, data interface{}) error {
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, fileData, 0644)
}
