package apidata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// ApiData fetches data from the API based on the provided endpoint and unmarshals it into the target structure.
func ApiData(endpoint string, target any, wg *sync.WaitGroup, errCh chan int) {
	defer wg.Done() // Mark this goroutine as done once the function exits.

	apiURL := "http://localhost:3000/api/" // Base URL of the API.
	fullURL := apiURL + endpoint

	// Make the GET request to the API endpoint.
	response, err := http.Get(apiURL + endpoint)
	if err != nil {
		fmt.Printf("API request error for endpoint %s: %s\n", fullURL, err.Error())
		if strings.Contains(err.Error(), "connection refused") {
			errCh <- http.StatusInternalServerError // Send internal server error code to the error channel.
		} else {
			errCh <- http.StatusBadRequest // Send bad request error code to the error channel.
		}
		return
	}
	defer response.Body.Close() // Ensure that the response body is closed.

	// Check the content type to verify that it's JSON.
	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		fmt.Printf("Unexpected content type %s for endpoint %s\n", contentType, fullURL)
		errCh <- http.StatusInternalServerError // Send internal server error code to the error channel.
		return
	}

	// Read the response body.
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("API reading error for endpoint %s: %s\n", fullURL, err.Error())
		errCh <- http.StatusInternalServerError // Send internal server error code to the error channel.
		return
	}

	// Unmarshal the JSON response data into the target structure.
	err = json.Unmarshal(responseData, target)
	if err != nil {
		fmt.Printf("API unmarshal error for endpoint %s: %s\n", fullURL, err.Error())
		fmt.Printf("Response body: %s\n", string(responseData)) // Print the response body for debugging.
		errCh <- http.StatusInternalServerError                 // Send internal server error code to the error channel.
		return
	}
}
