package handler

import (
	"encoding/json"
	"go-boot-category-api/model"
	"go-boot-category-api/service"
	"net/http"
	"strconv"
	"strings"
)

type categoryHandler struct {
	service service.Category
}

func NewCategoryHandler(service service.Category) *categoryHandler {
	return &categoryHandler{service: service}
}

// ValidateCategory validates the category data
func (h *categoryHandler) ValidateCategory(category *model.Category, isUpdate bool) string {
	if category.Name == "" {
		return "Category name is required"
	}
	if len(category.Name) < 3 {
		return "Category name must be at least 3 characters"
	}
	if len(category.Name) > 255 {
		return "Category name must not exceed 100 characters"
	}
	if len(category.Description) > 500 {
		return "Category description must not exceed 500 characters"
	}
	return ""
}

// HandleCategorys - GET /api/categories
func (h *categoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *categoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	Categorys, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Categorys)
}

func (h *categoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate category data
	if validationErr := h.ValidateCategory(&category, false); validationErr != "" {
		http.Error(w, validationErr, http.StatusBadRequest)
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// HandleCategoryByID - GET/PUT/DELETE /api/categories/{id}
func (h *categoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetByID - GET /api/categories/{id}
func (h *categoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	Category, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Category)
}

func (h *categoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var Category model.Category
	err = json.NewDecoder(r.Body).Decode(&Category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate category data
	if validationErr := h.ValidateCategory(&Category, true); validationErr != "" {
		http.Error(w, validationErr, http.StatusBadRequest)
		return
	}

	Category.ID = id
	err = h.service.Update(&Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Category)
}

// Delete - DELETE /api/categories/{id}
func (h *categoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Category deleted successfully",
	})
}
