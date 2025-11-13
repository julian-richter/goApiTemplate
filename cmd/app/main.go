package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/julian-richter/ApiTemplate/internal/config"
	"github.com/julian-richter/ApiTemplate/internal/db"
	model "github.com/julian-richter/ApiTemplate/internal/models/logentry"
	repo "github.com/julian-richter/ApiTemplate/internal/repos/logentry"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	pgPool, err := db.NewPostgresPool(cfg)
	if err != nil {
		log.Fatalf("Failed to open Postgres pool: %v", err)
	}
	defer pgPool.Close()

	valkeyCli, err := db.NewValkeyClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect Valkey: %v", err)
	}
	defer valkeyCli.Close()

	logRepo := repo.NewRepo(pgPool, repo.WithCache(valkeyCli, "app:"))

	app := fiber.New()

	// Search endpoint
	app.Get("/logs/search", func(c *fiber.Ctx) error {
		level := c.Query("level", "")
		messageContains := c.Query("message_contains", "")
		sinceStr := c.Query("since", "")
		untilStr := c.Query("until", "")
		limit := c.QueryInt("limit", 100)
		offset := c.QueryInt("offset", 0)

		var since, until *time.Time
		if sinceStr != "" {
			t, err := time.Parse(time.RFC3339, sinceStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "invalid since timestamp format",
					"details": sinceStr,
				})
			}
			since = &t
		}
		if untilStr != "" {
			t, err := time.Parse(time.RFC3339, untilStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "invalid until timestamp format",
					"details": untilStr,
				})
			}
			until = &t
		}

		params := repo.SearchParams{
			Level:           level,
			MessageContains: messageContains,
			Since:           since,
			Until:           until,
			Limit:           limit,
			Offset:          offset,
		}

		entries, err := logRepo.Search(context.Background(), params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "search failed",
				"details": err.Error(),
			})
		}

		// Make sure slice is non-nil
		if entries == nil {
			entries = []*model.LogEntry{}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data":   entries,
			"count":  len(entries),
			"offset": offset,
			"limit":  limit,
		})
	})

	app.Get("/logs/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid id")
		}
		entry, err := logRepo.GetByID(context.Background(), id, true, 5*time.Minute)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return c.JSON(entry)
	})

	app.Post("/logs", func(c *fiber.Ctx) error {
		var input model.LogEntry
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid body")
		}
		if input.Timestamp.IsZero() {
			input.Timestamp = time.Now()
		}
		if err := logRepo.Save(context.Background(), &input); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(input)
	})

	// New: Search endpoint
	app.Get("/logs/search", func(c *fiber.Ctx) error {
		level := c.Query("level", "")
		messageContains := c.Query("message_contains", "")
		sinceStr := c.Query("since", "")
		untilStr := c.Query("until", "")
		limit := c.QueryInt("limit", 100)
		offset := c.QueryInt("offset", 0)

		var since, until *time.Time
		if sinceStr != "" {
			t, err := time.Parse(time.RFC3339, sinceStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("invalid since timestamp format")
			}
			since = &t
		}
		if untilStr != "" {
			t, err := time.Parse(time.RFC3339, untilStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("invalid until timestamp format")
			}
			until = &t
		}

		params := repo.SearchParams{
			Level:           level,
			MessageContains: messageContains,
			Since:           since,
			Until:           until,
			Limit:           limit,
			Offset:          offset,
		}

		entries, err := logRepo.Search(context.Background(), params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(entries)
	})

	fmt.Printf("Server listening on port %s\n", cfg.App.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.App.Port)))
}
