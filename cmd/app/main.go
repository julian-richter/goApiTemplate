package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/julian-richter/ApiTemplate/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	// Create a new fiber app
	app := fiber.New()

	// Start the Fiber/v2 server
	log.Fatal(app.Listen(strings.TrimSpace(fmt.Sprintf(":%s", cfg.App.Port))))
}
