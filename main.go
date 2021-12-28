package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var recipes []Recipe
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

func init() {
	recipes = make([]Recipe, 0)
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
}

type Recipe struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Tags        []string           `json:"tags" bson:"tags"`
	Ingredients []string           `json:"ingredients" bson:"ingredients"`
	Instrctions []string           `json:"instructions" bson:"instructions"`
	PublishedAt time.Time          `json:"publishedAt" bson:"publishedAt"`
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while inserting new recipe"})
		return
	}
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// func UpdateRecipeHandler(c *gin.Context) {
// 	id := c.Param("id")
// 	var recipe Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}

// 	index := -1
// 	for i := 0; i < len(recipes); i++ {
// 		if recipes[i].ID == id {
// 			index = i
// 		}
// 	}

// 	if index == -1 {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Recipe Not Found"})
// 		return
// 	}
// 	recipe.ID = recipes[index].ID
// 	recipes[index] = recipe
// 	c.JSON(http.StatusOK, recipe)
// }

// func DeleteRecipeHandler(c *gin.Context) {

// 	id := c.Param("id")
// 	var recipe Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}

// 	index := -1
// 	for i := 0; i < len(recipes); i++ {
// 		if recipes[i].ID == id {
// 			index = i
// 		}
// 	}

// 	if index == -1 {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "recipe not found"})
// 	}

// 	recipes = append(recipes[:index], recipes[index+1:]...)
// 	c.JSON(http.StatusOK, gin.H{"message": "recipe has been deleted"})

// }

// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	listOfRecipes := make([]Recipe, 0)

// 	for i := 0; i < len(recipes); i++ {
// 		found := false
// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 			}
// 		}
// 		if found {
// 			listOfRecipes = append(listOfRecipes, recipes[i])
// 		}
// 	}

// 	c.JSON(http.StatusOK, listOfRecipes)
// }

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	// router.PUT("/recipes/:id", UpdateRecipeHandler)
	// router.DELETE("/recipes/:id", DeleteRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
