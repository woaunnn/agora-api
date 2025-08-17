package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func connectDB() (*gorm.DB, error) {
	username := "root"
	password := "55140239"
	host := "127.0.0.1"
	port := 3306
	database := "agoradb"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	log.Println("✅ Connected to MySQL via GORM!")
	return db, nil
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/api/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello, world!"})
	})

	app.Get("/api/user", func(c *fiber.Ctx) error {
		var users []User
		if result := db.Find(&users); result.Error != nil {
			return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		return c.JSON(users)
	})

	app.Post("/api/user", func(c *fiber.Ctx) error {
		name := c.Query("user")
		ageStr := c.Query("age")

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid age"})
		}

		user := User{Name: name, Age: age}
		if result := db.Create(&user); result.Error != nil {
			return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "User created successfully",
			"user":    user,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
