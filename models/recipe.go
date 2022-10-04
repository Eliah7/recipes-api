package models

import (
	"time"
)

// type Recipe struct {
// 	//swagger:ignore
// 	ID           primitive.ObjectID `json:"id" bson:"_id"`
// 	Name         string             `json:"name" bson:"name"`
// 	Tags         []string           `json:"tags" bson:"tags"`
// 	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
// 	Instructions []string           `json:"instructions" bson:"instructions"`
// 	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
// }

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}
