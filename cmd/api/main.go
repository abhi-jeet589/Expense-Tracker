package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abhi-jeet589/Expense-Tracker/internal/models"
	"github.com/abhi-jeet589/Expense-Tracker/internal/routes"
	"github.com/abhi-jeet589/Expense-Tracker/internal/webassets"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() *gorm.DB {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("error getting sql database handle: ", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		log.Fatal("error running migrations: ", err)
	}

	return db
}

func loadTemplates() *template.Template {
	tmpl := template.Must(template.New("").ParseFS(
		webassets.FS,
		"templates/*.tmpl",
		"templates/transactions/*.tmpl",
	))
	return tmpl
}

func DbInjector(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func main() {
	db := openDB()
	staticFS, err := fs.Sub(webassets.FS, "static")
	if err != nil {
		log.Fatal("error loading static assets: ", err)
	}

	router := gin.Default()
	router.Use(DbInjector(db))
	router.SetHTMLTemplate(loadTemplates())
	router.StaticFS("/static", http.FS(staticFS))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"healthy": true,
		})
	})

	routes.RegisterAPIRoutes(router.Group("/api"))
	routes.RegisterWebRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("error starting server: ", err)
	}
}
