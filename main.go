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
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// Load .env file untuk development
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_CONN", "host=localhost port=5432 user=postgres password=postgres dbname=mydb sslmode=disable")

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// DEBUG: Print config (hapus di production)
	log.Printf("Config - Port: %s", config.Port)
	log.Printf("Config - DB_CONN: %s", maskDBConn(config.DBConn))

	// Setup database dengan environment variable
	// Set DB_CONN ke environment variable agar database.InitDB() bisa baca
	if config.DBConn != "" {
		os.Setenv("DB_CONN", config.DBConn)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repository.NewProduct(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	categoryRepo := repository.NewCategory(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Category API",
			"docs":    "/api/products, /api/categories",
		})
	})

	log.Printf("Server running on port %s", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Helper function untuk mask password di logs
func maskDBConn(conn string) string {
	if len(conn) < 20 {
		return "***"
	}
	return conn[:20] + "..."
}
