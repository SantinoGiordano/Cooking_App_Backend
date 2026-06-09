package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func SearchRecipes(c *fiber.Ctx) error {
	query := c.Query("q")
	cuisine := c.Query("cuisine")
	difficulty := c.Query("difficulty")
	timeFilter := c.Query("time")

	collection := db.Collection("recipes")

	filter := bson.M{}

	// Search by recipe name
	if query != "" {
		filter["name"] = bson.M{
			"$regex":   query,
			"$options": "i",
		}
	}

	// Cuisine filter
	if cuisine != "" {
		filter["cuisine"] = cuisine
	}

	// Difficulty filter
	if difficulty != "" {
		diffInt, err := strconv.Atoi(difficulty)
		if err == nil {
			filter["difficulty"] = diffInt
		}
	}

	// Time filter
	switch timeFilter {
	case "short":
		filter["prepTime"] = bson.M{
			"$lte": 15,
		}

	case "medium":
		filter["prepTime"] = bson.M{
			"$gte": 20,
			"$lte": 40,
		}

	case "long":
		filter["prepTime"] = bson.M{
			"$gt": 40,
		}
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	defer cursor.Close(context.Background())

	var recipes []Recipe

	if err := cursor.All(context.Background(), &recipes); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if recipes == nil {
		recipes = []Recipe{}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(recipes),
		"data":    recipes,
	})
}