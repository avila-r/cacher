package main

import (
	"log"
	"os"

	"github.com/avila-r/xgo/pkg/api"
	"github.com/avila-r/xgo/pkg/cache"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Data struct {
	Data string `json:"data"`
	Id   int    `json:"id"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file can't be load")
	}

	options, _ := redis.ParseURL("redis://localhost:6379")

	cacher := cache.NewClient[Data](options)

	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
	})

	app.Post("/data", func(c *fiber.Ctx) error {
		var (
			data Data
		)

		_ = c.BodyParser(&data)

		cacher.Cache(cache.Insert[Data]{
			Key:  data.Data,
			Data: data,
		})

		return c.JSON(fiber.Map{
			"message": "success at saving data",
		})
	})

	app.Get("/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")

		data, _ := cacher.Uncache(cache.Query{
			Key: key,
		})

		return c.JSON(data)
	})

	url := os.Getenv("SERVER_URL")

	app.Listen(url)
}
