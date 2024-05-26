package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
		"model": "llava:latest",
		"prompt": "What is in this picture?",
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

	// Read and print the response
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(responseData))
}
