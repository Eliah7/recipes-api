//go:generate swagger generate spec
// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building- Distributed-Applications-in-Gin.
//
//  Schemes: http
//  Host: localhost:8080
//  BasePath: /
//  Version: 1.0.0
//  Contact: Eliah Mbwilo
// <eliahmbwilo7@gmail.com> https://eliah7.github.io
//
// Consumes:
//  - application/json
//
// Produces:
//  - application/json
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


var ctx context.Context
var err error
var client *mongo.Client

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Mongo Db")

	// get recipes from file and dump to db
	// file, _ := ioutil.ReadFile("recipes.json")
	// _ = json.Unmarshal([]byte(file), &recipes)

	// var listOfRecipes []interface{}

	// for _, recipe := range recipes {
	// 	listOfRecipes = append(listOfRecipes, recipe)
	// }

	// collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	// insertManyResults, err := collection.InsertMany(ctx, listOfRecipes)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Inserted recipes: ", len(insertManyResults.InsertedIDs))
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)

	c.JSON(http.StatusOK, recipe)
}

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes,
				recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})

}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
// '200':
//         description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		if err := cur.Decode(&recipe); err != nil{
			recipes = append(recipes, recipe)
		}
		
	}

	c.JSON(http.StatusOK, recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
