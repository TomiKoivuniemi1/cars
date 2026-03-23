package handlers

import (
	"cars/pkg/apidata"
	"cars/pkg/assets"
	"cars/pkg/cartypes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"text/template"
)

// CarHandler handles requests to the car details page.
func CarHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the URL path is correct and the method is GET
	if r.URL.Path != "/car" {
		assets.HandleError(w, http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		assets.HandleError(w, http.StatusMethodNotAllowed)
		return
	}

	// Parse the query from the URL to extract the car ID
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		assets.HandleError(w, http.StatusBadRequest)
		return
	}
	ModelIDStr := query.Get("id")
	if ModelIDStr == "" {
		assets.HandleError(w, http.StatusBadRequest)
		return
	}

	// Convert ID to integer
	ModelID, err := strconv.Atoi(ModelIDStr)
	if err != nil {
		assets.HandleError(w, http.StatusBadRequest)
		return
	}

	// Concurrently fetch data for the specific car model, manufacturers, and categories
	var model cartypes.Model
	var manufacturers []cartypes.Manufacturer
	var categories []cartypes.Category

	var wg sync.WaitGroup
	errCh := make(chan int, 3)
	wg.Add(3)
	go apidata.ApiData("models/"+strconv.Itoa(ModelID), &model, &wg, errCh)
	go apidata.ApiData("manufacturers", &manufacturers, &wg, errCh)
	go apidata.ApiData("categories", &categories, &wg, errCh)
	wg.Wait()

	close(errCh)
	if assets.ApiErrorFound(errCh) {
		assets.HandleError(w, <-errCh)
		return
	}

	// Populate the model with manufacturer and category details
	if model.ID != 0 {
		model.Manufacturer = manufacturers[model.ManufacturerID-1].Name
		model.Country = manufacturers[model.ManufacturerID-1].Country
		model.FoundingYear = manufacturers[model.ManufacturerID-1].FoundingYear
		model.Category = categories[model.CategoryID-1].Name
	}

	// Parse and execute the template
	html, err := template.ParseFiles("html/car.html")
	if err != nil {
		fmt.Println("Error parsing the html:", err)
		assets.HandleError(w, http.StatusInternalServerError)
		return
	}

	err = html.Execute(w, model)
	if err != nil {
		fmt.Println("Error executing the template:", err)
		assets.HandleError(w, http.StatusInternalServerError)
		return
	}
}
