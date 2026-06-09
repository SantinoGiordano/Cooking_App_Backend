package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// Recipe struct - adjust fields based on your actual collection structure
type Recipe struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name"`
	Ingredients []string `json:"ingredients" bson:"ingredients"`
	Instructions string   `json:"instructions" bson:"instructions"`
	PrepTime    int      `json:"prepTime" bson:"prepTime"`
	CookTime    int      `json:"cookTime" bson:"cookTime"`
	Difficulty  int   `json:"difficulty" bson:"difficulty"`
	Cuisine     string   `json:"cuisine" bson:"cuisine"`
	// Add more fields as needed, or use bson.M for flexible structure
}

// getAllRecipes fetches all recipes from the recipes collection
func getAllRecipes(c *fiber.Ctx) error {
	// Get the recipes collection
	collection := db.Collection("recipes")

	// Find all documents
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("❌ Error finding recipes:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch recipes",
		})
	}
	defer cursor.Close(context.Background())

	// Decode all recipes
	var recipes []Recipe
	if err = cursor.All(context.Background(), &recipes); err != nil {
		log.Println("❌ Error decoding recipes:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse recipes",
		})
	}

	// Return empty array if no recipes found
	if recipes == nil {
		recipes = []Recipe{}
	}

	log.Printf("✅ Found %d recipes", len(recipes))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    recipes,
		"count":   len(recipes),
	})
}

// Alternative: If you don't know the exact structure, use bson.M
func getAllRecipesFlexible(c *fiber.Ctx) error {
	collection := db.Collection("recipes")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch recipes",
		})
	}
	defer cursor.Close(context.Background())

	var recipes []bson.M
	if err = cursor.All(context.Background(), &recipes); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse recipes",
		})
	}

	if recipes == nil {
		recipes = []bson.M{}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    recipes,
		"count":   len(recipes),
	})
}