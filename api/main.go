package main
import (
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"urlshortener/routes"
)
func setupRoutes(app *fiber.App) {
	// Define your routes here
	app.Post("/api/shorten", routes.ShortenURL)
	app.Get("/:url", routes.ResolveURL)
}
func main() {
	err := godotenv.Load()
	if err!=nil {
		fmt.Println(err)
	}
	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
