package main

import (
	"context"
	"errors"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pgPool, err := db.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to open Postgres pool: %v", err)
	}
	defer pgPool.Close()

	// ------------------------------------------------------------
	// OPTIONAL VALKEY CACHE
	// ------------------------------------------------------------
	var logRepo *repo.Repo

	valkeyCli, err := db.NewValkeyClient(cfg)
	if err != nil {
		log.Printf("[warning] Valkey cache disabled: %v", err)
		logRepo = repo.NewRepo(pgPool) // no cache
	} else {
		defer valkeyCli.Close()
		log.Printf("[info] Valkey cache enabled")
		logRepo = repo.NewRepo(pgPool, repo.WithCache(valkeyCli, "app:"))
	}

	// ------------------------------------------------------------
	// HTTP SERVER
	// ------------------------------------------------------------
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
		ServerHeader:          "ApiTemplate",
	})

	// Search endpoint
	app.Get("/logs/search", func(c *fiber.Ctx) error {
		maxLimit := 500
		level := c.Query("level", "")
		messageContains := c.Query("message_contains", "")
		sinceStr := c.Query("since", "")
		untilStr := c.Query("until", "")
		limit := c.QueryInt("limit", 100)
		offset := c.QueryInt("offset", 0)

		if limit <= 0 {
			limit = 100
		}
		if limit > maxLimit {
			limit = maxLimit
		}

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

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()
		entries, err := logRepo.Search(ctx, params)
		if err != nil {
			log.Printf("search error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "search failed",
			})
		}

		// Proper "no results" error
		if len(entries) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "no log entries found",
				"details": params,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data":   entries,
			"count":  len(entries),
			"offset": offset,
			"limit":  limit,
		})
	})

	// Get a single log entry
	app.Get("/logs/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid id")
		}

		// use request context, not Background()
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		entry, err := logRepo.GetByID(ctx, id, true, 5*time.Minute)
		if err != nil {
			// distinguish "not found" from real errors
			if errors.Is(err, repo.ErrNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "log entry not found",
				})
			}

			log.Printf("get log by id error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to get log entry",
			})
		}

		return c.JSON(entry)
	})

	// Create log entry
	app.Post("/logs", func(c *fiber.Ctx) error {
		var input model.LogEntry
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid body")
		}

		if input.Timestamp.IsZero() {
			input.Timestamp = time.Now().UTC()
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		if err := logRepo.Save(ctx, &input); err != nil {
			log.Printf("save log entry error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to save log entry",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(input)
	})

	fmt.Printf("Server listening on port %s\n", cfg.App.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.App.Port)))
}
