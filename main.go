package main

import (
	"encoding/json"
	"go-boot-category-api/database"
	"go-boot-category-api/framework/handler"
	"go-boot-category-api/framework/repository"
	"go-boot-category-api/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env untuk development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, using environment variables")
	}

	// Setup database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Setup repositories, services, handlers
	productRepo := repository.NewProduct(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	categoryRepo := repository.NewCategory(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Setup router dengan middleware
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/", productHandler.HandleProductByID)
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Health check untuk Zeabur
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"database":  "connected",
		})
	})

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service":   "Category API",
			"version":   "1.0",
			"endpoints": []string{"/api/products", "/api/categories", "/health"},
		})
	})

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Add logging middleware
	handlerWithLogging := loggingMiddleware(mux)

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📡 Access URL: http://localhost:%s", port)
	log.Printf("🏥 Health check: http://localhost:%s/health", port)

	// Start server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handlerWithLogging,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("❌ Server failed:", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
