package triviafacts

import (
	"bytes"
	"chatbot/pkg/realtime"
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Get_Articles() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT public.get_articles()`).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_articles")

	return result, nil
}

func Get_Trivia() (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT public.get_trivia()`).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "get_trivia")

	return result, nil
}

func Get_FeatureImage(params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Pass whole struct to Postgres JSONB function
	if err := db.Raw(
		`SELECT public.get_triviaorfatcs_feature_image(?)`,
		params,
	).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Convert string fields to JSON map (if your sharedfunctions require it)
	sharedfunctions.ConvertStringToJSONMap(result)

	// Extract the function result
	result = sharedfunctions.GetMap(result, "get_triviaorfatcs_feature_image")

	return result, nil
}

// func Update_Articles(params *EditTriviaAndArticles) (map[string]any, error) {
func GetFeatureNameByID(id int) (string, error) {
	db := database.DB
	var featureName string
	err := db.Raw("SELECT feature_name FROM secondary_features WHERE id = ?", id).Scan(&featureName).Error
	return featureName, err
}

func Update_ArticlesOrTrivia(params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Get original featurename before update
	originalFeature, err := GetFeatureNameByID(params.ID)
	if err != nil {
		return nil, err
	}

	if err := db.Raw(`SELECT public.update_secondary_feature(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "update_secondary_feature")

	// Determine which hubs to update based on featurename changes
	updatedFeature := params.Featurename
	if originalFeature == updatedFeature {
		// Featurename didn't change - update only the relevant hub
		switch updatedFeature {
		case "Articles":
			dataArticle, err := Get_Articles()
			if err != nil {
				return nil, err
			}
			realtime.ArticlesHub.Publish(dataArticle)
		case "Trivia Facts":
			dataTrivia, err := Get_Trivia()
			if err != nil {
				return nil, err
			}
			realtime.TriviaHub.Publish(dataTrivia)
		}
	} else {
		// Featurename changed - update both hubs
		dataArticle, err := Get_Articles()
		if err != nil {
			return nil, err
		}
		dataTrivia, err := Get_Trivia()
		if err != nil {
			return nil, err
		}
		realtime.ArticlesHub.Publish(dataArticle)
		realtime.TriviaHub.Publish(dataTrivia)
	}

	return result, nil
}

func Delete_triviaorfacts(params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any
	// Delete via Postgres function
	if err := db.Raw("SELECT public.delete_triviaorfatcs(?)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "delete_triviaorfatcs")
	if err := broadcasting(params.Featurename); err != nil {
		return nil, err
	}
	return result, nil
}

// inserting articel function
func Insert_ArticleOrTrivia(params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Insert into Postgres via function
	if err := db.Raw(`SELECT public.insert_article(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "insert_article")

	if err := broadcasting(params.Featurename); err != nil {
		return nil, err
	}

	return result, nil
}

func broadcasting(featureName string) error {
	var data map[string]any
	var err error

	switch featureName {
	case "Articles":
		data, err = Get_Articles()
		realtime.ArticlesHub.Publish(data)
	case "Trivia Facts":
		data, err = Get_Trivia()
		realtime.TriviaHub.Publish(data)
	default:
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func GetGitHubRepo(id int64) (string, string, string, string, error) {
	db := database.DB
	var result map[string]any

	// Call the Postgres function with integer ID
	if err := db.Raw(`SELECT baseurl.get_github_repo(?)`, id).Scan(&result).Error; err != nil {
		return "", "", "", "", err
	}

	// Convert string fields to JSON map if needed
	sharedfunctions.ConvertStringToJSONMap(result)

	// Extract the actual JSONB result from the Postgres function
	repoData := sharedfunctions.GetMap(result, "get_github_repo")

	// Extract the required fields
	githubToken, _ := repoData["github_token"].(string)
	owner, _ := repoData["owner"].(string)
	repo, _ := repoData["repo"].(string)
	path, _ := repoData["path"].(string)

	return githubToken, owner, repo, path, nil
}

func UploadToGitHub(filename string, content []byte) (string, error) {
	// Fetch GitHub repo details from DB function
	gitHubUrl, err := sharedfunctions.GetBaseUrl(11)
	if err != nil {
		return "", err
	}

	gitHubImageLink, err := sharedfunctions.GetBaseUrl(14)
	if err != nil {
		return "", err
	}

	gitHubToken, err := sharedfunctions.GetBasicAuth(2)
	if err != nil {
		return "", err
	}

	// Build the path and URL
	url := gitHubUrl + filename

	// Prepare request body
	body := map[string]string{
		"message": "Upload " + filename,
		"content": base64.StdEncoding.EncodeToString(content),
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))

	req.Header.Set("Authorization", "Bearer "+gitHubToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("GitHub API error: %s", string(respBytes))
	}

	// Construct raw file URL
	rawURL := gitHubImageLink + filename
	return rawURL, nil
}

func DeleteFromGitHub(part1 string) error {

	gitHubUrl, err := sharedfunctions.GetBaseUrl(11)
	if err != nil {
		return err
	}

	gitHubToken, err := sharedfunctions.GetBasicAuth(2)
	if err != nil {
		return err
	}

	url := gitHubUrl + part1

	// GET file SHA
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+gitHubToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub GET error: %s", string(body))
	}

	var fileInfo struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
		return err
	}

	// DELETE file
	body := map[string]string{
		"message": "Delete " + url,
		"sha":     fileInfo.SHA,
	}
	jsonBody, _ := json.Marshal(body)

	reqDel, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
	reqDel.Header.Set("Authorization", "Bearer "+gitHubToken)
	reqDel.Header.Set("Accept", "application/vnd.github+json")

	respDel, err := client.Do(reqDel)
	if err != nil {
		return err
	}
	defer respDel.Body.Close()

	respBody, _ := io.ReadAll(respDel.Body)

	if respDel.StatusCode >= 400 {
		return fmt.Errorf("GitHub DELETE error: %s", string(respBody))
	}

	return nil
}
func UpdateFileOnGitHub(filename string, content []byte) (string, error) {

	gitHubUrl, err := sharedfunctions.GetBaseUrl(11)
	if err != nil {
		return "", err
	}

	gitHubImageLink, err := sharedfunctions.GetBaseUrl(14)
	if err != nil {
		return "", err
	}

	gitHubToken, err := sharedfunctions.GetBasicAuth(2)
	if err != nil {
		return "", err
	}

	url := gitHubUrl + filename
	client := &http.Client{}

	var sha string

	// Try GET file SHA
	reqGet, _ := http.NewRequest("GET", url, nil)
	reqGet.Header.Set("Authorization", "Bearer "+gitHubToken)
	reqGet.Header.Set("Accept", "application/vnd.github+json")
	respGet, err := client.Do(reqGet)
	if err != nil {
		return "", err
	}
	defer respGet.Body.Close()

	if respGet.StatusCode == 200 {
		var fileInfo struct {
			SHA string `json:"sha"`
		}
		if err := json.NewDecoder(respGet.Body).Decode(&fileInfo); err == nil {
			sha = fileInfo.SHA
		}
	} else if respGet.StatusCode != 404 {
		// Only fail if it's not a "Not Found"
		body, _ := io.ReadAll(respGet.Body)
		return "", fmt.Errorf("GitHub GET error: %s", string(body))
	}

	// Now PUT (create or update)
	body := map[string]string{
		"message": "Upload/Update " + filename,
		"content": base64.StdEncoding.EncodeToString(content),
	}
	if sha != "" {
		body["sha"] = sha
	}
	jsonBody, _ := json.Marshal(body)

	reqPut, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	reqPut.Header.Set("Authorization", "Bearer "+gitHubToken)
	reqPut.Header.Set("Accept", "application/vnd.github+json")

	respPut, err := client.Do(reqPut)
	if err != nil {
		return "", err
	}
	defer respPut.Body.Close()

	if respPut.StatusCode >= 400 {
		respBytes, _ := io.ReadAll(respPut.Body)
		return "", fmt.Errorf("GitHub UPDATE error: %s", string(respBytes))
	}

	rawURL := gitHubImageLink + filename
	return rawURL, nil
}
