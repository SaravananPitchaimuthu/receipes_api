package main

import (
	"context"
	"log"
	"os"

	"github.com/SaravananPitchaimuthu/receipes_api/receipes_api/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

var recipeHandler *handlers.RecipeHandler

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	recipeHandler = handlers.NewRecipeHandler(collection, ctx)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipeHandler.NewRecipeHandler)
	router.GET("/recipes", recipeHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)
	router.Run()
}
