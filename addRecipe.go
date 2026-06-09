package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func AddRecipe(c *fiber.Ctx) error {
	var recipe Recipe

	if err := c.BodyParser(&recipe); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	collection := db.Collection("recipes")

	result, err := collection.InsertOne(context.Background(), recipe)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to insert recipe",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      result.InsertedID,
	})
}
