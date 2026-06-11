package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var client *mongo.Client
var db *mongo.Database

func connectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using Render environment variables")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("❌ MONGO_URI environment variable is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Mongo connection error:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Mongo ping failed:", err)
	}

	log.Println("✅ Connected to MongoDB")

	db = client.Database("cookingData")
}

func main() {
	connectDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,OPTIONS,PUT,DELETE",
		// AllowCredentials: true,
	}))
	app.Get("/api/recipes", getAllRecipes)
	app.Get("/api/recipes/search", SearchRecipes)
	app.Post("/api/recipes/by-ingredients", SearchRecipesByIngredients)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
