package assets

import (
	"fmt"
	"net/http"
)

// HandleError sends a plain text error response to the client based on the given HTTP status code.
func HandleError(w http.ResponseWriter, status int) {
	if status <= 0 {
		status = http.StatusInternalServerError
	}

	var ErrMsg struct {
		Code    int
		Message string
	}
	ErrMsg.Code = status

	switch status {
	case http.StatusNotFound:
		ErrMsg.Message = "Oops! The page you're looking for doesn't exist."
	case http.StatusMethodNotAllowed:
		ErrMsg.Message = "This method is not allowed for the requested resource."
	case http.StatusBadRequest:
		ErrMsg.Message = "The request could not be understood by the server."
	case http.StatusInternalServerError:
		ErrMsg.Message = "We're experiencing technical difficulties. Please try again later."
	default:
		ErrMsg.Message = "An unexpected error occurred. Please contact support."
	}

	// Write the header once, with the correct status code
	w.WriteHeader(status)

	// Set the content type to plain text
	w.Header().Set("Content-Type", "text/plain")

	// Write the error message to the response
	_, err := fmt.Fprintf(w, "Error %d: %s", ErrMsg.Code, ErrMsg.Message)
	if err != nil {
		fmt.Println("Failed to write error message:", err)
	}
}

// ApiErrorFound checks for errors
func ApiErrorFound(errCh <-chan int) bool {
	for errCode := range errCh {
		if errCode != 0 {
			return true
		}
	}
	return false
}
