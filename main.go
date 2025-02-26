package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var (
	waitingRoomCap = 10
	waitingRoomKey = "waiting_room"
	waitingRoomTTL = 1 * time.Minute
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to Redis on localhost:6379")

	return client
}

func NewWaitingRoomMiddleware(redisClient *redis.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		uniqId := c.Get("Authorization")
		if uniqId == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Get the number of users in the waiting room
		waitingRoomCount, err := redisClient.HLen(context.Background(), waitingRoomKey).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Check if the user is already in the waiting room
		for _, v := range redisClient.HVals(context.Background(), waitingRoomKey).Val() {
			if v == uniqId {
				return c.Next()
			}
		}

		// Check if waiting room is available
		if int(waitingRoomCount) >= waitingRoomCap {
			return c.Status(fiber.StatusServiceUnavailable).SendString("Waiting room is full")
		}

		// Add user to the waiting room
		err = redisClient.HSet(context.Background(), waitingRoomKey, uniqId, uniqId).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Set TTL for the user in the waiting room
		redisClient.HExpire(context.Background(), waitingRoomKey, waitingRoomTTL, uniqId)

		return c.Next()
	}
}

func main() {
	app := fiber.New()
	redisClient := NewRedisClient()

	app.Use(NewWaitingRoomMiddleware(redisClient))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	go func() {
		time.Sleep(5 * time.Millisecond)
		log.Println(`Starting server on localhost:3000`)
	}()
	log.Fatal(app.Listen(":3000"))
}
