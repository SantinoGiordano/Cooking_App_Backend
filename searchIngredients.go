package main 

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var RecipeCollection *mongo.Collection

type IngredientSearchRequest struct {
	Ingredients []string `json:"ingredients"`
}

func SearchRecipesByIngredients(c *fiber.Ctx) error {
	var req IngredientSearchRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if len(req.Ingredients) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "At least one ingredient is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Match recipes that contain ANY of the ingredients
	filter := bson.M{
		"ingredients": bson.M{
			"$in": req.Ingredients,
		},
	}

	cursor, err := RecipeCollection.Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var recipes []Recipe

	if err := cursor.All(ctx, &recipes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(recipes),
		"data":    recipes,
	})
}