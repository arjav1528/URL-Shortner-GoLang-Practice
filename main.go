package main

import (
	"log"
	"os"
	"time"

	model "github.com/arjav1528/URL-Shortner-GoLang-Practice/src/models"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	_ = godotenv.Load()
	app.Get("/", getRoot)
	app.Post("/shorten", addURL)
	app.Get("/:hash", redirectURL)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

func getRoot(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello",
	})
}

func addURL(c fiber.Ctx) error {
	var url model.URL

	if err := c.Bind().Body(&url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Body",
		})
	}

	hash, err := url.AddURL()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"hash":   hash,
		"url":    url.URL,
		"expiry": url.Expiry.Local().String(),
	})
}

func redirectURL(c fiber.Ctx) error {
	var u model.URL

	hash := c.Params("hash")
	if hash == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing hash",
		})
	}

	url, err := u.GetModel(hash)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if time.Now().After(url.Expiry) {
		return c.Status(fiber.ErrConflict.Code).JSON(fiber.Map{
			"error": "Shortner Expired",
		})
	} else {
		return c.Redirect().Status(fiber.StatusFound).To(url.URL)
	}
}
