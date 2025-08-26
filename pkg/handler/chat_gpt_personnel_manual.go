// package handler

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"sync"

// 	"github.com/PullRequestInc/go-gpt3"
// 	"github.com/gofiber/fiber/v2"
// )

// type Response struct {
// 	QuestionOptions []string
// 	Answer          string
// }

// var (
// 	responses        map[string]Response
// 	conversation     []string
// 	conversationLock sync.Mutex
// )
// var (
// 	// Add a flag to track if the first conversation has occurred
// 	firstConversationCompleted bool
// )

// func init() {
// 	if err := library(); err != nil {
// 		log.Fatal("Failed to initialize responses: ", err)
// 	}
// }

// func library() error {
// 	// Read the responses from the file
// 	var err error
// 	responses, err = readResponsesFromFile("responses.txt")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // HTTP handler for the chat endpoint
// func ChatHandler(c *fiber.Ctx) error {
// 	if c.Method() != "POST" {
// 		return c.Status(http.StatusMethodNotAllowed).SendString("Invalid request method")
// 	}

// 	userInput := c.FormValue("input")
// 	response, statusCode := getResponse(userInput)
// 	return c.Status(statusCode).SendString(response)
// }

// // Read responses from a file
// func readResponsesFromFile(filename string) (map[string]Response, error) {
// 	// Read the responses from the file
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	responses := make(map[string]Response)
// 	scanner := bufio.NewScanner(file)
// 	var currentKey string
// 	var currentResponse Response
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
// 			if currentKey != "" {
// 				responses[currentKey] = currentResponse
// 			}
// 			currentKey = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
// 			currentResponse = Response{}
// 			options := strings.Split(currentKey, "/")
// 			currentResponse.QuestionOptions = options
// 		} else {
// 			currentResponse.Answer += line + "\n"
// 		}
// 	}
// 	if currentKey != "" {
// 		responses[currentKey] = currentResponse
// 	}

// 	return responses, scanner.Err()
// }

// func getResponse(userInput string) (string, int) {
// 	var matchingResponses []string

// 	// Check if there are matching responses in responses.txt
// 	matchingResponseFound := false
// 	for _, response := range responses {
// 		for _, option := range response.QuestionOptions {
// 			if strings.ToLower(userInput) == strings.ToLower(option) {
// 				matchingResponses = append(matchingResponses, response.Answer)
// 				matchingResponseFound = true
// 			}
// 		}
// 	}

// 	// If matching responses are found in responses.txt, combine and return them
// 	if matchingResponseFound {
// 		startLine := "Thank you for your message.\n\n"
// 		allResponses := strings.Join(matchingResponses, "")
// 		localResponse := startLine + allResponses
// 		return localResponse, http.StatusOK
// 	}

// 	// If no matching response is found in responses.txt, use ChatGPT
// 	apiKey := os.Getenv("API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("Missing API KEY")
// 	}

// 	client := gpt3.NewClient(apiKey)
// 	ctx := context.Background()

// 	// Use conversation as context
// 	conversationLock.Lock()
// 	conversationContext := strings.Join(conversation, "\n")
// 	conversationContext += "\nUser: " + userInput
// 	conversationLock.Unlock()

// 	// Append "CA-GABAY" to the conversation
// 	conversationContext += "\nCA-GABAY: "

// 	gptResponse, err := fetchChatGptResponse(client, ctx, conversationContext)
// 	if err != nil {
// 		log.Printf("Error: %s\n", err)
// 		return "", http.StatusInternalServerError
// 	}

// 	// Remove "Bot:" prefix from the ChatGPT response
// 	if strings.HasPrefix(gptResponse, "Bot:") {
// 		gptResponse = strings.TrimPrefix(gptResponse, "Bot:")
// 	}

// 	// Append the user input and ChatGPT response to the conversation
// 	conversationLock.Lock()
// 	conversation = append(conversation, ""+userInput)
// 	conversation = append(conversation, "CA-GABAY"+gptResponse)
// 	conversationLock.Unlock()

// 	return gptResponse, http.StatusOK
// }

// func fetchChatGptResponse(client gpt3.Client, ctx context.Context, input string) (string, error) {
// 	targetTokenCount := 200 // Set the desired maximum token count

// 	completion, err := client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
// 		Prompt:      []string{input},
// 		MaxTokens:   gpt3.IntPtr(targetTokenCount),
// 		Temperature: gpt3.Float32Ptr(0),
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}

// 	responseText := completion.Choices[0].Text

// 	// Check if the response exceeds the target token count
// 	tokens := strings.Split(responseText, " ")
// 	if len(tokens) > targetTokenCount {
// 		// If the response exceeds the target token count, truncate it to the desired length
// 		responseText = strings.Join(tokens[:targetTokenCount], " ")
// 	}

// 	// Check if this is the first conversation
// 	if !firstConversationCompleted {
// 		// Append the additional message only in the first conversation
// 		additionalMessage := "If you are seeking information within the Personnel Policies and Procedures Manual, you may enter 'Table of Contents' to access the document's index or check your message if you are mistaken."
// 		responseText += " " + additionalMessage

// 		// Set the flag to indicate that the first conversation has been completed
// 		firstConversationCompleted = true
// 	}

// 	return responseText, nil
// }

// // Function to check if farewell words are present in the response
// func containsFarewellWords(response string) bool {
// 	farewellWords := []string{"goodbye", "farewell", "bye"}
// 	for _, word := range farewellWords {
// 		if strings.Contains(response, word) {
// 			return true
// 		}
// 	}
// 	return false
// }

// Chronological Order based on file
package handler

import (
	"bufio"
	"chatbot/pkg/logs"
	"chatbot/pkg/models/errors"
	"chatbot/pkg/models/response"
	"chatbot/pkg/models/status"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	//"time"

	"unicode"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

type ChatRequest struct {
	StaffID string `json:"staffid"`
	Input   string `json:"input"`
}

type Response struct {
	QuestionOptions []string
	Answer          string
}

var (
	responses []Response
	//conversation      []string
	conversationLock  sync.Mutex
	userConversations = make(map[string][]string)
)

var (
	// Add a flag to track if the first conversation has occurred
	// firstConversationCompleted bool
	conversationResetInterval = 5 * time.Minute
)

func init() {
	if err := personnelManual(); err != nil {
		log.Fatal("Failed to initialize responses: ", err)
	}
	userConversations = make(map[string][]string) // Initialize userConversations map
	go resetConversationHistory()
}

func personnelManual() error {
	// Read the responses from the file
	var err error
	responses, err = readResponsesFromFile("pkg/personnel_manual/responses.txt")
	if err != nil {
		return err
	}
	return nil
}

// Function to reset the conversation history periodically
func resetConversationHistory() {
	for {
		// // Sleep for the specified interval
		time.Sleep(conversationResetInterval)

		// Reset the conversation history
		conversationLock.Lock()
		for user, conversation := range userConversations {
			// Check if there are more than 5 entries
			if len(conversation) > 5 {
				// Create a new slice to store the last 5 entries
				newConversation := conversation[len(conversation)-5:]
				userConversations[user] = newConversation
			}
		}
		conversationLock.Unlock()
	}
}

// HTTP handler for the chat endpoint
func ChatHandler(c *fiber.Ctx) error {
	chatRequest := new(ChatRequest)
	if err := c.BodyParser(chatRequest); err != nil {
		fmt.Println("Failed to parse chat request", err.Error())
		return c.Status(401).JSON(response.ResponseModel{
			RetCode: "401",
			Message: status.RetCode401,
			Data: errors.ErrorModel{
				Message:   "Failed to parse chat request",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	staffID := c.Params("id")
	responses, err := getResponse( /*chatRequest.StaffID,*/ chatRequest.Input)
	if err != nil {
		logs.LOSLogs(c, "Chat", staffID, "500", err.Error())
		return c.Status(500).JSON(response.ResponseModel{
			RetCode: "500",
			Message: status.RetCode500,
			Data: errors.ErrorModel{
				Message:   "Failed to generate chatGPT response.",
				IsSuccess: false,
				Error:     err,
			},
		})
	}

	logs.LOSLogs(c, "Chat", staffID, "200", (chatRequest.Input + "\n" + responses))
	return c.Status(200).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success!",
		Data:    responses,
	})
}

// Read responses from a file into a slice
func readResponsesFromFile(filename string) ([]Response, error) {
	// Read the responses from the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var responses []Response
	var currentResponse Response
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if currentResponse.Answer != "" {
				responses = append(responses, currentResponse)
			}
			currentResponse = Response{}
			options := strings.Split(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"), "/")
			currentResponse.QuestionOptions = options
		} else {
			currentResponse.Answer += line + "\n"
		}
	}
	if currentResponse.Answer != "" {
		responses = append(responses, currentResponse)
	}

	return responses, scanner.Err()
}

func getResponse( /*userID string,*/ userInput string) (string, error) {
	var matchingResponses []string

	// Check if there are matching responses in responses.txt
	matchingResponseFound := false
	for _, response := range responses {
		for _, option := range response.QuestionOptions {
			if func() string {
				isASCII, hasUpper := true, false
				for i := 0; i < len(userInput); i++ {
					c := userInput[i]
					if c >= utf8.RuneSelf {
						isASCII = false
						break
					}
					hasUpper = hasUpper || ('A' <= c && c <= 'Z')
				}
				if isASCII {
					if !hasUpper {
						return userInput
					}
					var (
						b   strings.Builder
						pos int
					)
					b.Grow(len(userInput))
					for i := 0; i < len(userInput); i++ {
						c := userInput[i]
						if 'A' <= c && c <= 'Z' {
							c += 'a' - 'A'
							if pos < i {
								b.WriteString(userInput[pos:i])
							}
							b.WriteByte(c)
							pos = i + 1
						}
					}
					if pos < len(userInput) {
						b.WriteString(userInput[pos:])
					}
					return b.String()
				}
				return strings.Map(unicode.ToLower, userInput)
			}() == strings.ToLower(option) {
				matchingResponses = append(matchingResponses, response.Answer)
				matchingResponseFound = true
			}
		}
	}

	// If matching responses are found in responses.txt, combine and return them
	if !matchingResponseFound {
		startLine := "No policy found for your input. Please select titles from the Table of Contents."
		allResponses := strings.Join(matchingResponses, "")
		localResponse := startLine + allResponses
		return localResponse, nil
	}

	// // If no matching response is found in responses.txt, use ChatGPT
	// apiKey := utils.GetEnv("API_KEY")
	// if apiKey == "" {
	// 	log.Fatal("Missing API KEY")
	// }

	// client, err := gpt35.NewClient(apiKey)
	// if err != nil {
	// 	log.Fatalf("Failed to create GPT-3.5 Turbo client: %v", err)
	// 	return "", err
	// }

	// // Use conversation as context
	// conversationLock.Lock()
	// conversationContext := strings.Join(userConversations[userID], "\n")
	// conversationLock.Unlock()

	// // Ensure "CA-GABAY" is part of the conversation context
	// if !strings.Contains(conversationContext, "CA-GABAY") {
	// 	conversationContext += "\nCA-GABAY: "
	// }

	// conversationContext += "\nUser: " + userInput

	// // Append "CA-GABAY" to the conversation
	// conversationLock.Lock()
	// conversationContext += "\nCA-GABAY: "
	// conversationLock.Unlock()

	// gptResponse, err := fetchChatGptResponse(client, conversationContext)
	// if err != nil {
	// 	log.Printf("Error: %s\n", err)
	// 	return "", err
	// }

	// // Remove "Bot:" prefix from the ChatGPT response
	// if func() bool {
	// 	var prefix string = "Bot:"
	// 	return len(gptResponse) >= len(prefix) && gptResponse[0:len(prefix)] == prefix
	// }() {
	// 	gptResponse = strings.TrimPrefix(gptResponse, "Bot:")
	// }

	// // Append the user input and ChatGPT response to the user's conversation
	// conversationLock.Lock()
	// userConversations[userID] = append(userConversations[userID], userInput)
	// userConversations[userID] = append(userConversations[userID], "CA-GABAY "+gptResponse)
	// conversationLock.Unlock()

	// return gptResponse, nil
	startLine := "Thank you for your message.\n\n"
	allResponses := strings.Join(matchingResponses, "")
	localResponse := startLine + allResponses
	return localResponse, nil
}

// func fetchChatGptResponse(client *gpt35.Client, input string) (string, error) {
// 	targetTokenCount := 500 // Set the desired maximum token count

// 	// Define a request for the chat completion
// 	request := &gpt35.Request{
// 		Model: gpt35.ModelGpt35Turbo,
// 		Messages: []*gpt35.Message{
// 			{
// 				Role:    gpt35.RoleUser,
// 				Content: input, // Use the input as the user's message
// 			},
// 		},
// 		MaxTokens:   gpt35.MaxTokensGpt35Turbo, // Set the maximum token count for the response
// 		Temperature: 0,
// 	}

// 	// Make a request to the ChatGPT API
// 	response, err := client.GetChat(request)
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	}

// 	// Get the response message from the API
// 	responseText := response.Choices[0].Message.Content

// 	// Check if the response exceeds the target token count
// 	tokens := strings.Split(responseText, " ")
// 	if len(tokens) > targetTokenCount {
// 		// If the response exceeds the target token count, truncate it to the desired length
// 		responseText = strings.Join(tokens[:targetTokenCount], " ")
// 	}

// 	// Check if this is the first conversation
// 	if !firstConversationCompleted {
// 		// Append the additional message only in the first conversation
// 		additionalMessage := "If you are seeking information within the Personnel Policies and Procedures Manual, you can send 'Table of Contents' to view the Titles."
// 		responseText += " " + additionalMessage

// 		// Set the flag to indicate that the first conversation has been completed
// 		firstConversationCompleted = true
// 	} else {
// 		// Check if the user's input mentions the manual or policies and procedures
// 		if strings.Contains(strings.ToLower(input), "Personnel Policies and Procedures Manual") {
// 			// If the user mentioned the manual, provide information about the table of contents
// 			tableOfContentsMessage := "You can send 'Table of Contents' to view the Personnel Policies and Procedures Manual's Index."
// 			responseText += " " + tableOfContentsMessage
// 		} else if strings.Contains(strings.ToLower(input), "CA-GABAY") {
// 			// If the user mentioned the manual, provide information about the table of contents
// 			tableOfContentsMessage := "CARD GABAY Application"
// 			responseText += " " + tableOfContentsMessage
// 		}
// 	}

// 	return responseText, nil
// }
