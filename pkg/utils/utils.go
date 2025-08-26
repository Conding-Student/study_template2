// Package utils provides ...
package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// func GetEnv(key string) string {
// 	err := godotenv.Load("project_env_files/.env")

// 	if err != nil {
// 		fmt.Println("Error loading .env file")
// 		log.Fatalf("Error loading .env file")
// 	}

// 	return os.Getenv(key)
// }

func GetEnv(key string) string {
	// Load variables from the main .env file
	err := godotenv.Load("project_env_files/.env")

	if err != nil {
		fmt.Println("Error loading project_env_files/.env file")
		log.Fatalf("Error loading project_env_files/.env file")
	}

	// Get the value of ENVIRONMENT variable
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		fmt.Println("ENVIRONMENT variable not set")
		log.Fatalf("ENVIRONMENT variable not set")
	}

	// Load variables from the corresponding .env file based on ENVIRONMENT
	subEnvFile := fmt.Sprintf("project_env_files/.env-%s", env)
	err = godotenv.Load(subEnvFile)

	if err != nil {
		fmt.Printf("Error loading %s file\n", subEnvFile)
		log.Fatalf("Error loading %s file\n", subEnvFile)
	}

	return os.Getenv(key)
}
