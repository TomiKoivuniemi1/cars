package assets

import (
	"cars/pkg/cartypes"
	"net/url"
	"strconv"
)

// Filter filters a list of car models based on query parameters such as manufacturer and category.
// It returns a list of models that satisfy all the provided filter criteria.
func Filter(query url.Values, models []cartypes.Model) []cartypes.Model {
	var filteredResult []cartypes.Model

	for _, model := range models {
		if !matchManufacturer(model, query["manufacturerID"]) {
			continue
		}
		if !matchCategory(model, query["categoryID"]) {
			continue
		}
		filteredResult = append(filteredResult, model)
	}

	return filteredResult
}

// matchManufacturer determines if a model's manufacturer matches any of the specified manufacturer IDs in the query.
// Returns true if no manufacturer IDs are specified, implying no filtering on manufacturer.
func matchManufacturer(model cartypes.Model, ManufacturerID []string) bool {
	if len(ManufacturerID) == 0 || ManufacturerID[0] == "" {
		return true
	}
	return ManufacturerID[0] == strconv.Itoa(model.ManufacturerID)
}

// matchCategory determines if a model's category matches any of the specified category IDs in the query.
// Returns true if no category IDs are specified, implying no filtering on category.
func matchCategory(model cartypes.Model, CategoryID []string) bool {
	if len(CategoryID) == 0 || CategoryID[0] == "" {
		return true
	}
	return CategoryID[0] == strconv.Itoa(model.CategoryID)
}
