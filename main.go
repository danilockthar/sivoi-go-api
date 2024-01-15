package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model        // embedes basic ID int , CreatedAt, UpdatedAt
	Title      string `json:"title" binding:"required" required:"El titulo no debe estar vacío"`
	Status     string `json:"status" binding:"required"`
}

/* Validator input */
var (
	g = galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	customizer = g.Validator(Todo{})
)

func main() {
	dsn := "host=db user=myuser password=mypassword dbname=myapp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&Todo{})

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Dani!",
		})
	})

	r.GET("/todos", func(c *gin.Context) {
		var todos []Todo

		// Query all todos from the database
		if err := db.Order("created_at desc").Find(&todos).Error; err != nil {
			c.JSON(500, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		// Respond with the list of todos
		c.JSON(200, todos)
	})

	r.POST("/todos", func(c *gin.Context) {
		var input Todo

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"status":  "failed_binding",
				"message": customizer.DecryptErrors(err),
			})
			return
		}

		// Create the record in the database
		if err := db.Create(&input).Error; err != nil {
			fmt.Println("Error creating todo:", err)
			c.JSON(400, gin.H{
				"success": false,
				"status":  "cannot_create_db",
				"message": "Cannot create todo",
			})
			return
		}

		// Respond with success
		c.JSON(200, gin.H{
			"success": true,
			"message": "Todo created successfully",
			"todo":    input,
		})
	})

	r.Run()
}
