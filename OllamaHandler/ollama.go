package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var ConversationHistoryMap map[string][]string = make(map[string][]string)
var CreateModel_url string
var GenerateResponse_url string

func getFileContentAsString(filename string) (string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	// Read the content of the file
	content := make([]byte, fileSize)
	_, err = file.Read(content)
	if err != nil {
		return "", err
	}

	// Convert content to string
	fileContent := string(content)
	return fileContent, nil
}

func makeOllamaRequest(data map[string]interface{}, url string) map[string]interface{} {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return map[string]interface{}{}
	}

	fmt.Println("Sending Request to Ollama")

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return map[string]interface{}{}
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Receiving Response from Ollama")
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return map[string]interface{}{}
		}

		var responseData map[string]interface{}
		err = json.Unmarshal(responseBody, &responseData)
		if err != nil {
			fmt.Println("Error unmarshalling response JSON:", err)
			return map[string]interface{}{}
		}

		return responseData
	} else {
		fmt.Println("Error:", response.Status)
		return map[string]interface{}{}
	}
}

func relativeToAbsolutePath(relativePath string) string {
	// Get the absolute path
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}

	// Check if the file or directory exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		fmt.Println("File or directory does not exist.")
	} else if err != nil {
		fmt.Println("Error checking file or directory existence:", err)
	}

	return absolutePath
}

func CreateModels(models map[string]string) {
	for key, value := range models {
		fmt.Printf("Creating %s model\n", key)
		fmt.Println(relativeToAbsolutePath(value))

		// Get the content of the file as a string
		content, err := getFileContentAsString(relativeToAbsolutePath(value))
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		data := map[string]interface{}{
			"name":      key,
			"modelfile": content,
			"stream":    false,
		}

		responseData := makeOllamaRequest(data, CreateModel_url)
		if len(responseData) == 0 {
			fmt.Println("Received Empty Response")
			continue
		}

		actualResponse := responseData["status"].(string)

		fmt.Printf("Created %s %s\n", key, actualResponse)
	}
}

func GenerateResponse(prompt string, model string) string {
	conversationHistory := ConversationHistoryMap[model]
	conversationHistory = append(conversationHistory, prompt)

	fullPrompt := ""
	for _, entry := range conversationHistory {
		fullPrompt += entry + "\n"
	}

	data := map[string]interface{}{
		"model":  model,
		"stream": false,
		"prompt": fullPrompt,
	}

	responseData := makeOllamaRequest(data, GenerateResponse_url)
	if len(responseData) == 0 {
		return ""
	}

	actualResponse := responseData["response"].(string)
	conversationHistory = append(conversationHistory, actualResponse)
	if len(conversationHistory) >= 200 {
		fmt.Printf("Pruning conversation history for %v", model)
		conversationHistory = conversationHistory[len(conversationHistory)/2:]
	}
	ConversationHistoryMap[model] = conversationHistory
	return actualResponse
}
