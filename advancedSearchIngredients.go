package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func AdvancedSearchRecipesByIngredient(c *fiber.Ctx) error {
	ingredient := c.Query("ingredient")
	cuisine := c.Query("cuisine")
	difficulty := c.Query("difficulty")
	timeFilter := c.Query("time")
	matchAll := c.Query("matchAll") == "true" // If true, must contain ALL ingredients

	collection := db.Collection("recipes")

	filter := bson.M{}

	// Search by ingredient(s)
	if ingredient != "" {
		// Split by comma to support multiple ingredients
		ingredients := strings.Split(ingredient, ",")
		
		// Trim whitespace from each ingredient
		for i := range ingredients {
			ingredients[i] = strings.TrimSpace(ingredients[i])
		}

		if matchAll {
			// Must contain ALL ingredients
			var conditions []bson.M
			for _, ing := range ingredients {
				if ing != "" {
					conditions = append(conditions, bson.M{
						"ingredients": bson.M{
							"$regex":   ing,
							"$options": "i",
						},
					})
				}
			}
			if len(conditions) > 0 {
				filter["$and"] = conditions
			}
		} else {
			// Must contain AT LEAST ONE ingredient
			var conditions []bson.M
			for _, ing := range ingredients {
				if ing != "" {
					conditions = append(conditions, bson.M{
						"ingredients": bson.M{
							"$regex":   ing,
							"$options": "i",
						},
					})
				}
			}
			if len(conditions) > 0 {
				filter["$or"] = conditions
			}
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