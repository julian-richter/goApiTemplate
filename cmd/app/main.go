package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/julian-richter/ApiTemplate/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()
	// Create a new fiber app
	app := fiber.New()

	// Start the Fiber/v2 server
	log.Fatal(app.Listen(":" + strconv.Itoa(cfg.App.Port)))
}
