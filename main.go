package main

import (
	"bucket-token-algorthm-golang/limitter"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()

	capacity := 10
	rate := 1
	limiter := limitter.NewTokenBucket(capacity, rate)

	bucketTokenLimiterMiddleware := func(c *fiber.Ctx) error {
		if !limiter.TakeTokens(1) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message":     "Too many requests",
				"description": fmt.Sprintf("Your max request number is %d. Please try again in %d seconds later", capacity, rate),
			})
		}
		return c.Next()
	}

	app.Use(bucketTokenLimiterMiddleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hai")
	})

	log.Println("Server running on localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
