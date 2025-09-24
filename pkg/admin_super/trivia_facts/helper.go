package triviafacts

import (
	"bytes"
	"chatbot/pkg/realtime"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GitHubConfig struct {
	BaseURL   string
	ImageLink string
	Token     string
}

func getGitHubConfig() (*GitHubConfig, error) {
	baseURL, err := sharedfunctions.GetBaseUrl(11)
	if err != nil {
		return nil, err
	}

	imageLink, err := sharedfunctions.GetBaseUrl(14)
	if err != nil {
		return nil, err
	}

	token, err := sharedfunctions.GetBasicAuth(2)
	if err != nil {
		return nil, err
	}

	return &GitHubConfig{
		BaseURL:   baseURL,
		ImageLink: imageLink,
		Token:     token,
	}, nil
}

func executeGitHubRequest(method, url, token string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

func getFileSHA(url, token string) (string, error) {
	respBytes, err := executeGitHubRequest("GET", url, token, nil)
	if err != nil {
		return "", err
	}

	var fileInfo struct {
		SHA string `json:"sha"`
	}
	if err := json.Unmarshal(respBytes, &fileInfo); err != nil {
		return "", err
	}

	return fileInfo.SHA, nil
}
func broadcasting(id string, featureName string) error {
	var data map[string]any
	var err error

	switch featureName {
	case "Articles":
		data, err = Get_Articles()
		realtime.MainHub.Publish(id, "get_articles", data)

	case "Trivia Facts":
		data, err = Get_Trivia()
		realtime.MainHub.Publish(id, "get_trivia", data)
	default:
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func Get_FeatureImage(params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Pass whole struct to Postgres JSONB function
	if err := db.Raw(`SELECT public.get_triviaorfatcs_feature_image(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert string fields to JSON map (if your sharedfunctions require it)
	sharedfunctions.ConvertStringToJSONMap(result)

	// Extract the function result
	result = sharedfunctions.GetMap(result, "get_triviaorfatcs_feature_image")

	return result, nil
}
