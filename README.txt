package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

)



func main() {
	// Read the responses from the file
	responses, err := readResponsesFromFile("responses.txt")
	if err != nil {
		log.Fatal("Failed to read responses from file: ", err)
	}

	// Start the chat loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nuser: ")
		scanner.Scan()
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Chatbot: Goodbye!")
			break
		}
		if strings.ToLower(userInput) == "quit" {
			fmt.Println("Chatbot: Goodbye! Next time don't quit. You can do that! :D ")
			break
		}
		if strings.ToLower(userInput) == "goodbye" {
			fmt.Println("Chatbot: Goodbye too! come back again next time :D ")
			break
		}

		response := getResponse(userInput, responses)
		fmt.Println("\nChatbot:", response)
	}
}

// Read responses from a file
func readResponsesFromFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	responses := make(map[string]string)
	scanner := bufio.NewScanner(file)
	var currentKey string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentKey = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
		} else {
			responses[currentKey] += line + "\n"
		}
	}

	return responses, scanner.Err()
}

// Get a response based on user input
func getResponse(userInput string, responses map[string]string) string {
	for question, answer := range responses {
		if strings.ToLower(userInput) == strings.ToLower(question) {
			return answer
		}
	}

	return "Paumanhin! Pakisuri ang iyong word/ words, baka mayroong mali sa mga letra, o, magtungo muli sa tables of contents \n para makuha ang exact word na dapat ay ilagay sa input. Salamat."
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var responses map[string]string

func main() {
	// Read the responses from the file
	var err error
	responses, err = readResponsesFromFile("responses.txt")
	if err != nil {
		log.Fatal("Failed to read responses from file: ", err)
	}

	// Set up HTTP routes
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/", defaultHandler)

	// Start the server
	port := ":8000" // I-update ang port base sa iyong preference
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// HTTP handler for the chat endpoint
func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userInput := r.FormValue("input")
	response := getResponse(userInput)
	fmt.Fprintf(w, response)
}

// HTTP handler for the default endpoint
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the chatbot server!")
}

// Read responses from a file
func readResponsesFromFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	responses := make(map[string]string)
	scanner := bufio.NewScanner(file)
	var currentKey string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentKey = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
		} else {
			responses[currentKey] += line + "\n"
		}
	}

	return responses, scanner.Err()
}

// Get a response based on user input
func getResponse(userInput string) string {
	for question, answer := range responses {
		if strings.ToLower(userInput) == strings.ToLower(question) {
			return answer
		}
	}

	return "Paumanhin! Pakisuri ang iyong word/words, baka mayroong mali sa mga letra, o, magtungo muli sa table of contents para makuha ang eksaktong salita na dapat ilagay sa input. Salamat."
}


API: "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", // Your OpenAI API key
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var responses map[string]string

func main() {
	// Read the responses from the file
	var err error
	responses, err = readResponsesFromFile("responses.txt")
	if err != nil {
		log.Fatal("Failed to read responses from file: ", err)
	}

	// Set up HTTP routes
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/", defaultHandler)

	// Start the server
	port := ":8000" // Update the port based on your preference
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// HTTP handler for the chat endpoint
func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userInput := r.FormValue("input")
	response := getResponse(userInput)
	fmt.Fprintf(w, response)
}

// HTTP handler for the default endpoint
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the chatbot server!")
}

// Read responses from a file
func readResponsesFromFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	responses := make(map[string]string)
	scanner := bufio.NewScanner(file)
	var currentKey string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentKey = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
		} else {
			responses[currentKey] += line + "\n"
		}
	}

	return responses, scanner.Err()
}

// Get a response based on user input
func getResponse(userInput string) string {
	for question, answer := range responses {
		if strings.ToLower(userInput) == strings.ToLower(question) {
			return answer
		}
	}

	// Use OpenAI GPT or other fallback method to generate response
	openai := &OpenAI{
		API: "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", // Your OpenAI API key
	}
	response, err := openai.ChatGPT(userInput)
	if err != nil {
		log.Println("Failed to generate response using OpenAI GPT:", err)
		return "An error occurred while generating the response."
	}

	return response
}

// OpenAI struct and chatGPT function for OpenAI GPT integration
type OpenAI struct {
	API string
}

type GPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func (o *OpenAI) ChatGPT(prompt string) (string, error) {
	apiURL := "https://api.openai.com/v1/completions"

	payload := map[string]any{
		"model":             "text-davinci-003",
		"prompt":            prompt,
		"temperature":       0,
		"max_tokens":        2000,
		"top_p":             1,
		"frequency_penalty": 0.0,
		"presence_penalty":  0.0,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.API)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response GPTResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Text, nil
	}

	return "", nil

}
