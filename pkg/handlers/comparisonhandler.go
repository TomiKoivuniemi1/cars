package handlers

import (
	"cars/pkg/apidata"
	"cars/pkg/assets"
	"cars/pkg/cartypes"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"text/template"
)

// ComparisonHandler handles requests to the compare page.
func ComparisonHandler(w http.ResponseWriter, r *http.Request) {
	// Check for path and method
	if r.URL.Path != "/compare" {
		assets.HandleError(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		assets.HandleError(w, http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	// Parse the query from the URL to extract model IDs
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		assets.HandleError(w, http.StatusBadRequest)
		return
	}

	// Get the model IDs from the query parameters
	modelIDs := query["IDs"]
	if len(modelIDs) == 0 {
		pageData := cartypes.Cartypes{
			Path: "compare",
		}
		html, _ := template.ParseFiles("./html/compare.html")
		html.Execute(w, pageData)
		return
	}

	// Fetch data concurrently
	var models []cartypes.Model
	var manufacturers []cartypes.Manufacturer
	var categories []cartypes.Category

	var wg sync.WaitGroup
	errCh := make(chan int, 2)
	wg.Add(2)
	go apidata.ApiData("manufacturers", &manufacturers, &wg, errCh)
	go apidata.ApiData("categories", &categories, &wg, errCh)
	wg.Wait()
	close(errCh)

	// Check for errors from concurrent operations
	if assets.ApiErrorFound(errCh) {
		assets.HandleError(w, http.StatusInternalServerError)
		return
	}

	// Fetch the models based on IDs using the existing ModelId function
	models = apidata.ModelId(modelIDs, models, manufacturers, categories, w)

	// Populate the page data with the fetched data and models for comparison
	pageData := cartypes.Cartypes{
		Path:          "compare",
		Categories:    categories,
		Manufacturers: manufacturers,
		Models:        models,
	}

	// Parse and execute the comparison template
	html, err := template.ParseFiles("./html/compare.html")
	if err != nil {
		fmt.Println("Error parsing the HTML template:", err)
		assets.HandleError(w, http.StatusInternalServerError)
		return
	}

	err = html.Execute(w, pageData)
	if err != nil {
		fmt.Println("Error executing the HTML template:", err)
		assets.HandleError(w, http.StatusInternalServerError)
		return
	}
}
