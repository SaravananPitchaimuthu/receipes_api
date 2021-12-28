package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Receipe struct {
	Name        string    `json:"name"`
	Tags        []string  `json:"tags"`
	Ingredients []string  `json:"ingredients"`
	Instrctions []string  `json:"instructions"`
	PublishedAt time.Time `json:"publishedAt"`
}

func main() {
	router := gin.Default()
	router.Run()
}
