package core

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName    string `json:"app_name"`
	APIPrefix  string `json:"api_prefix"`
	APIVersion string `json:"api_version"`
	AppSecret  string `json:"app_secret"`
	DB         struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	TokenConfig struct {
		Duration int64  `json:"duration"` // Token duration in seconds
		Secret   string `json:"secret"`   // Secret key for signing tokens
	}
	AuthConfig struct {
		HeaderPrefix string `json:"header_prefix"` // Prefix for the authorization header
	}
	Mock struct {
		Enable   bool   `json:"enable"`    // Enable mock mode
		MockUser string `json:"mock_user"` // Mock user for testing
	}
}

func findFileRecursively(dir, filename string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		fullPath := dir + "/" + entry.Name()
		if entry.IsDir() {
			foundPath, err := findFileRecursively(fullPath, filename)
			if err == nil {
				return foundPath, nil
			}

			if !errors.Is(err, os.ErrNotExist) {
				return "", err
			}
		} else if entry.Name() == filename {
			return fullPath, nil
		}
	}
	return "", os.ErrNotExist
}

func findPathToEnvFile(filename string) string {
	pathToEnvFile := os.Getenv("ENV_FILE_PATH")
	if pathToEnvFile == "" {
		log.Println("ENV_FILE_PATH not set, using default path")
		dir, err := os.Getwd()
		if err != nil {
			panic("Failed to get current working directory: " + err.Error())
		}
		maybeRoot := strings.TrimSuffix(dir, "/backend/cmd")
		foundPath, err := findFileRecursively(maybeRoot, filename)
		if err != nil {
			panic("Failed to read directory: " + err.Error())
		}
		log.Println("FoundPath: ", foundPath)
		return foundPath
	}
	return pathToEnvFile
}
func InitEnv(filename string) {
	pathFound := findPathToEnvFile(filename)
	err := godotenv.Load(pathFound)
	if err != nil {
		panic("Failed to load environment variables: " + err.Error())
	}
	log.Printf("%s Starts!", os.Getenv("APP_NAME"))
}

func GetAppConfig() *AppConfig {

	cfg := &AppConfig{
		AppName:    os.Getenv("APP_NAME"),
		APIPrefix:  os.Getenv("API_PREFIX"),
		APIVersion: os.Getenv("API_VERSION"),
		AppSecret:  os.Getenv("APP_SECRET"),
		DB: struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			Name     string `json:"name"`
		}{
			Host:     "localhost",
			Port:     5432,
			User:     "dbuser",
			Password: "dbpassword",
			Name:     "myappdb",
		},
		TokenConfig: struct {
			Duration int64  `json:"duration"` // Token duration in seconds
			Secret   string `json:"secret"`   // Secret key for signing tokens,
		}{
			Duration: 604800, // 7 days in seconds
			Secret:   "mysecretkey",
		},
		AuthConfig: struct {
			HeaderPrefix string `json:"header_prefix"` // Prefix for the authorization header
		}{
			HeaderPrefix: "Bearer ",
		},
		Mock: struct {
			Enable   bool   `json:"enable"`    // Enable mock mode
			MockUser string `json:"mock_user"` // Mock user for testing
		}{
			Enable:   true,
			MockUser: "mockuser",
		},
	}
	return cfg
}
