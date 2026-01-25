package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ==================== Data Models ====================
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ==================== In-Memory Database ====================
var categoriesDB = []*Category{
	{ID: 1, Name: "Makanan", Description: "Category Makanan"},
	{ID: 2, Name: "Minuman", Description: "Category Minuman"},
}

// ==================== Helper Functions ====================
func validateCategory(category *Category) error {
	name := strings.TrimSpace(category.Name)

	if name == "" {
		return fmt.Errorf("Name is required")
	}
	if len(name) > 255 {
		return fmt.Errorf("Name cannot exceed 255 characters")
	}
	if len(category.Description) > 500 {
		return fmt.Errorf("Description cannot exceed 500 characters")
	}

	return nil
}

func findCategoryByID(id int) (*Category, int) {
	for index, category := range categoriesDB {
		if category.ID == id {
			return category, index
		}
	}
	return nil, -1
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	writeJSONResponse(w, statusCode, map[string]string{
		"error": message,
	})
}

// ==================== Category Handlers ====================
// GET /api/categories - Get all categories
func handleGetCategories(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, categoriesDB)
}

// POST /api/categories - Create new category
func handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category

	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validateCategory(&newCategory); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate new ID
	newCategory.ID = len(categoriesDB) + 1
	categoriesDB = append(categoriesDB, &newCategory)

	writeJSONResponse(w, http.StatusCreated, newCategory)
}

// GET /api/categories/{id} - Get category by ID
func handleGetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, _ := findCategoryByID(id)
	if category == nil {
		writeErrorResponse(w, http.StatusNotFound, "Category not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, category)
}

// PUT /api/categories/{id} - Update category
func handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var updatedData Category
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validateCategory(&updatedData); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	_, index := findCategoryByID(id)
	if index == -1 {
		writeErrorResponse(w, http.StatusNotFound, "Category not found")
		return
	}

	// Preserve the ID
	updatedData.ID = id
	categoriesDB[index] = &updatedData

	writeJSONResponse(w, http.StatusOK, updatedData)
}

// DELETE /api/categories/{id} - Delete category
func handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	_, index := findCategoryByID(id)
	if index == -1 {
		writeErrorResponse(w, http.StatusNotFound, "Category not found")
		return
	}

	// Remove from slice
	categoriesDB = append(categoriesDB[:index], categoriesDB[index+1:]...)

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Category successfully deleted",
	})
}

// ==================== Health Check Handler ====================
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "API is running",
	})
}

// ==================== Main Application ====================
func main() {
	mux := http.NewServeMux()

	// Register routes with Go 1.22 enhanced routing
	mux.HandleFunc("GET /api/categories", handleGetCategories)
	mux.HandleFunc("POST /api/categories", handleCreateCategory)
	mux.HandleFunc("GET /api/categories/{id}", handleGetCategoryByID)
	mux.HandleFunc("PUT /api/categories/{id}", handleUpdateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}", handleDeleteCategory)
	mux.HandleFunc("GET /health", handleHealthCheck)

	// Print server information
	fmt.Println("Server running on: http://localhost:8080")
	fmt.Println("\nAvailable Endpoints:")
	fmt.Println("  GET    /api/categories       - List all categories")
	fmt.Println("  POST   /api/categories       - Create new category")
	fmt.Println("  GET    /api/categories/{id}  - Get category by ID")
	fmt.Println("  PUT    /api/categories/{id}  - Update category")
	fmt.Println("  DELETE /api/categories/{id}  - Delete category")
	fmt.Println("  GET    /health               - Health check")

	// Start server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
