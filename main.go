package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

// In lieu of a go SDK for Ollama,
// we create a stand-in.
type ApiResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	Context            []int  `json:"context"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

func main() {
	// Command line flags
	imagePath := flag.String("image", "", "Path to the image file")
	baseURL := flag.String("baseurl", "http://localhost:11434", "Base URL of the API")

	flag.Parse()

	// Validate image path
	if *imagePath == "" {
		fmt.Println("Image path is required")
		os.Exit(1)
	}

	// Open the image file
	file, err := os.Open(*imagePath)
	if err != nil {
		fmt.Printf("Failed to open image file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the image file
	imageData, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read image file: %s\n", err)
		os.Exit(1)
	}

	// Encode the image to base64
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	// Prepare the API request
	url := *baseURL + "/api/generate"
	requestBody := fmt.Sprintf(`{
		"model": "llava",
		"prompt": "Provide an exhaustive description of the computer software image including identifying all objects and describing them and their relationships",
		"stream": false,
		"images": ["%s"]
	}`, encodedImage)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(requestBody))
	if err != nil {
		fmt.Printf("Failed to create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read and parse the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %s\n", err)
		os.Exit(1)
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(responseData, &apiResponse); err != nil {
		fmt.Printf("Failed to parse response JSON: %s\n", err)
		os.Exit(1)
	}

	// Print the response in a readable format
	fmt.Printf("Model: %s\n", apiResponse.Model)
	fmt.Printf("Created At: %s\n", apiResponse.CreatedAt)
	fmt.Printf("Response: %s\n", apiResponse.Response)
	fmt.Printf("Done: %t\n", apiResponse.Done)
	fmt.Printf("Done Reason: %s\n", apiResponse.DoneReason)
	fmt.Printf("Total Duration: %.1f s\n", float64(apiResponse.TotalDuration/1e9))
	fmt.Printf("Load Duration: %.1f s\n", float64(apiResponse.LoadDuration/1e9))
	fmt.Printf("Prompt Evaluation Count: %d\n", apiResponse.PromptEvalCount)
	fmt.Printf("Prompt Evaluation Duration: %.1f s\n", float64(apiResponse.PromptEvalDuration/1e9))
	fmt.Printf("Evaluation Count: %d\n", apiResponse.EvalCount)
	fmt.Printf("Evaluation Duration: %.1f s\n", float64(apiResponse.EvalDuration/1e9))
}
