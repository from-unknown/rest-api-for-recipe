package constants

import (
	"RestAPIForRecipe/models"
)

type GetRecipesResponse struct {
	Recipe []*RecipeWithoutTime `json:"recipe"`
}

type RecipeWithoutTime struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	MakingTime  string `json:"preparation_time"`
	Serves      string `json:"serves"`
	Ingredients string `json:"ingredients"`
	Cost        int    `json:"cost"`
}

type PatchRecipe struct {
	Title       string `json:"title"`
	MakingTime  string `json:"preparation_time"`
	Serves      string `json:"serves"`
	Ingredients string `json:"ingredients"`
	Cost        int    `json:"cost"`
}

type GetRecipeByIDResponse struct {
	Message string               `json:"message"`
	Recipe  []*RecipeWithoutTime `json:"recipe"`
}

type PostRecipesResponse struct {
	Message string           `json:"message"`
	Recipe  []*models.Recipe `json:"recipe"`
}

type PatchRecipesResponse struct {
	Message string         `json:"message"`
	Recipe  []*PatchRecipe `json:"recipe"`
}

type ErrorPostRecipesResponse struct {
	Message  string `json:"message"`
	Required string `json:"required"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

var GetRecipeByIDMessage = "Recipe details by id"

var RecipeCreateSuccess = "Recipe successfully created!"
var RecipeUpdateSuccess = "Recipe successfully updated!"

var ErrorParameterRequired = "title, preparation_time, serves, ingredients, cost"
var ErrorCreationFailed = "Recipe creation failed!"

var RecipeDeleteSuccess = "Recipe successfully removed!"
var RecipeNotFound = "No recipe found"
