package apidata

import (
	"cars/pkg/assets"
	"cars/pkg/cartypes"
	"net/http"
	"sync"
)

// ModelId concurrently fetches data for all the models specified by the given IDs.
// It populates the models slice with data, including category and manufacturer names,
// and writes errors to the provided HTTP response writer.
func ModelId(IDs []string, models []cartypes.Model, manufacturers []cartypes.Manufacturer, categories []cartypes.Category, w http.ResponseWriter) []cartypes.Model {
	var selectedModels []cartypes.Model // Slice to hold the selected models
	var wg sync.WaitGroup               // WaitGroup to manage concurrent API requests

	idCount := len(IDs)              // Number of IDs to process
	errCh := make(chan int, idCount) // Channel to capture errors from API requests

	// Check if the IDs slice is empty or contains a single empty string, return an empty slice if true
	if idCount == 1 && IDs[0] == "" {
		return selectedModels
	}

	// Pre-populate the resulting slice to avoid out of bounds errors
	for i := 0; i < idCount; i++ {
		selectedModels = append(selectedModels, cartypes.Model{})
	}

	// Add the number of IDs to the WaitGroup counter
	wg.Add(idCount)
	// Concurrently fetch data for each model ID from the API
	for i, id := range IDs {
		go ApiData("models/"+id, &selectedModels[i], &wg, errCh)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	close(errCh) // Close the error channel

	// Check if any API errors were encountered and handle them
	if assets.ApiErrorFound(errCh) {
		assets.HandleError(w, <-errCh) // Handle the first error found
		return selectedModels
	}

	// Process the fetched models: remove any with ID 0 and set category and manufacturer names
	for i := idCount - 1; i >= 0; i-- {
		if selectedModels[i].ID == 0 {
			// Remove models with ID 0 (indicating an error or missing data)
			selectedModels = append(selectedModels[:i], selectedModels[i+1:]...)
		} else {
			// Set the category and manufacturer names for the remaining models
			selectedModels[i].Category = categories[selectedModels[i].CategoryID-1].Name
			selectedModels[i].Manufacturer = manufacturers[selectedModels[i].ManufacturerID-1].Name
		}
	}
	return selectedModels
}
