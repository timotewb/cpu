package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func EnvVariables(fullPath string) {
    
    // Read the contents of the .env file
    content, err := os.ReadFile(filepath.Join(fullPath, ".env"))
    if err != nil {
        log.Fatalf("Error reading .env file: %v", err)
    }
    
    // Split the content into lines
    lines := strings.Split(string(content), "\n")
    
    // Parse each line
    for _, line := range lines {
        parts := strings.Fields(line)
        if len(parts) == 2 {
            key := parts[0]
            value := parts[1]
            
            // Remove quotes if present
            if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
                value = value[1 : len(value)-1]
            } else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
                value = value[1 : len(value)-1]
            }
			os.Setenv(key, value)
        }
    }
}