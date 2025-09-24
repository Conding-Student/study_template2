package triviafacts

import (
	"chatbot/pkg/sharedfunctions"
	"chatbot/pkg/utils/go-utils/database"
	"encoding/base64"
	"encoding/json"
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

func Update_ArticlesOrTrivia(satffid string, params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	if err := db.Raw(`SELECT public.update_secondary_feature(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}
	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "update_secondary_feature")

	if err := broadcasting(satffid, params.Featurename); err != nil {
		return nil, err
	}

	return result, nil
}

func Delete_triviaorfacts(staffid string, params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any
	// Delete via Postgres function
	if err := db.Raw("SELECT public.delete_triviaorfatcs(?)", params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "delete_triviaorfatcs")
	if err := broadcasting(staffid, params.Featurename); err != nil {
		return nil, err
	}
	return result, nil
}

// inserting articel function
func Insert_ArticleOrTrivia(staffid string, params *TriviaAndArticles) (map[string]any, error) {
	db := database.DB
	var result map[string]any

	// Insert into Postgres via function
	if err := db.Raw(`SELECT public.insert_article(?)`, params).Scan(&result).Error; err != nil {
		return nil, err
	}

	sharedfunctions.ConvertStringToJSONMap(result)
	result = sharedfunctions.GetMap(result, "insert_article")

	if err := broadcasting(staffid, params.Featurename); err != nil {
		return nil, err
	}

	return result, nil
}

func UploadToGitHub(filename string, content []byte) (string, error) {
	config, err := getGitHubConfig()
	if err != nil {
		return "", err
	}

	url := config.BaseURL + filename
	body := map[string]string{
		"message": "Upload " + filename,
		"content": base64.StdEncoding.EncodeToString(content),
	}
	jsonBody, _ := json.Marshal(body)

	_, err = executeGitHubRequest("PUT", url, config.Token, jsonBody)
	if err != nil {
		return "", err
	}

	return config.ImageLink + filename, nil
}

func DeleteFromGitHub(filename string) error {
	config, err := getGitHubConfig()
	if err != nil {
		return err
	}

	url := config.BaseURL + filename
	sha, err := getFileSHA(url, config.Token)
	if err != nil {
		return err
	}

	body := map[string]string{
		"message": "Delete " + filename,
		"sha":     sha,
	}
	jsonBody, _ := json.Marshal(body)

	_, err = executeGitHubRequest("DELETE", url, config.Token, jsonBody)
	return err
}

func UpdateFileOnGitHub(filename string, content []byte) (string, error) {
	config, err := getGitHubConfig()
	if err != nil {
		return "", err
	}

	url := config.BaseURL + filename
	body := map[string]string{
		"message": "Upload/Update " + filename,
		"content": base64.StdEncoding.EncodeToString(content),
	}

	// Try to get SHA for existing file
	if sha, err := getFileSHA(url, config.Token); err == nil {
		body["sha"] = sha
	}

	jsonBody, _ := json.Marshal(body)
	_, err = executeGitHubRequest("PUT", url, config.Token, jsonBody)
	if err != nil {
		return "", err
	}

	return config.ImageLink + filename, nil
}
