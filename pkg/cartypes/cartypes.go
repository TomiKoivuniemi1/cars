package cartypes

// Cartypes is a struct for the various types of cars needed to be handled in the the path of the car manufacturers, categories, and models.
type Cartypes struct {
	Path          string         // Path to the data source
	Manufacturers []Manufacturer // List of car manufacturers
	Categories    []Category     // List of car categories
	Models        []Model        // List of car models
}

// Manufacturer represents a car manufacturer.
type Manufacturer struct {
	ID           int    // Unique identifier for the manufacturer
	Name         string // Name of the manufacturer
	Country      string // Country where the manufacturer is based
	FoundingYear int    // Year the manufacturer was founded
}

// Category represents a category of cars.
type Category struct {
	ID   int    // Unique identifier for the category
	Name string // Name of the category
}

// Model represents a specific car model.
type Model struct {
	ID             int            // Unique identifier for the car model
	Name           string         // Name of the car model
	ManufacturerID int            // Identifier for the manufacturer
	Manufacturer   string         // Name of the manufacturer
	CategoryID     int            // Identifier for the category
	Category       string         // Name of the category
	Year           int            // Year the car model was released
	Specifications Specifications // Specifications of the car model
	Image          string         // URL to an image of the car model
	Country        string         // Country of the manufacturer
	FoundingYear   int            // Founding year of the manufacturer
}

// Specifications represents the detailed specifications of a car model.
type Specifications struct {
	Engine       string // Engine type of the car
	Horsepower   int    // Horsepower of the car
	Transmission string // Transmission type of the car
	Drivetrain   string // Drivetrain type of the car
}

// ErrMsg represents an error message with a code and a message.
type ErrMsg struct {
	Code    int    // Error code
	Message string // Error message
}
