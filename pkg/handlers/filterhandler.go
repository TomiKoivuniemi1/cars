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

// FilterHandler handles requests to the filter page. It fetches data for models, manufacturers,
// and categories, filters the models based on query parameters, and then renders the index
// template with the filtered data.
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the URL path is correct and the method is GET
	if r.URL.Path != "/filter" {
		assets.HandleError(w, http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		assets.HandleError(w, http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	// Parse the query from the URL
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		assets.HandleError(w, http.StatusBadRequest)
		return
	}

	// Parse the index template
	html, err := template.ParseFiles("./html/index.html")
	if err != nil {
		fmt.Println("Error parsing the html:", err)
		assets.HandleError(w, http.StatusInternalServerError)
		return
	}

	// Concurrently apidata data for models, manufacturers, and categories
	var models []cartypes.Model
	var manufacturers []cartypes.Manufacturer
	var categories []cartypes.Category

	var wg sync.WaitGroup
	errCh := make(chan int, 3)
	wg.Add(3)
	go apidata.ApiData("models", &models, &wg, errCh)
	go apidata.ApiData("manufacturers", &manufacturers, &wg, errCh)
	go apidata.ApiData("categories", &categories, &wg, errCh)
	wg.Wait()

	close(errCh)
	if assets.ApiErrorFound(errCh) {
		assets.HandleError(w, <-errCh)
		return
	}

	// Filter the models based on query parameters and populate the page data
	var pageData cartypes.Cartypes
	pageData.Path = "search"
	pageData.Categories = categories
	pageData.Manufacturers = manufacturers
	pageData.Models = assets.Filter(query, models)

	// Execute the template with the filtered page data
	err = html.Execute(w, pageData)
	if err != nil {
		fmt.Println("Error executing the template:", err)
		assets.HandleError(w, http.StatusInternalServerError)
		return
	}
}
