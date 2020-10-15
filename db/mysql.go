package db

import (
	"RestAPIForRecipe/constants"
	"RestAPIForRecipe/models"
	"database/sql"
	"errors"
	"os"
)

type SqlHandler struct {
	db *sql.DB
}

func NewSqlHandler() *SqlHandler {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err.Error())
	}

	return &SqlHandler{db: db}
}

func (sh *SqlHandler) GetRecipes() ([]*models.Recipe, error) {
	rows, err := sh.db.Query("SELECT id, title, making_time, serves, ingredients, cost, created_at, updated_at" +
		" FROM recipes")
	if err != nil {
		panic(err.Error())
	}

	var recipes []*models.Recipe
	for rows.Next() {
		recipe := models.Recipe{}
		err = rows.Scan(&recipe.ID, &recipe.Title, &recipe.MakingTime, &recipe.Serves, &recipe.Ingredients, &recipe.Cost,
			&recipe.CreatedAt, &recipe.UpdatedAt)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

func (sh *SqlHandler) GetRecipeByID(id int) ([]*models.Recipe, error) {
	rows, err := sh.db.Query("SELECT id, title, making_time, serves, ingredients, cost, created_at, updated_at "+
		"FROM recipes where id = ?", id)
	if err != nil {
		panic(err.Error())
	}

	var recipes []*models.Recipe
	for rows.Next() {
		recipe := models.Recipe{}
		err = rows.Scan(&recipe.ID, &recipe.Title, &recipe.MakingTime, &recipe.Serves, &recipe.Ingredients, &recipe.Cost,
			&recipe.CreatedAt, &recipe.UpdatedAt)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

func (sh *SqlHandler) InsertRecipe(title string, makingTime string, serves string, ingredients string, cost int) ([]*models.Recipe, error) {
	// Prepare stmt for inserting data
	stmt, err := sh.db.Prepare("INSERT INTO recipes (title, making_time, serves, ingredients, cost" +
		") VALUES( ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(title, makingTime, serves, ingredients, cost)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	recipe, err := sh.GetRecipeByID(int(id))
	return recipe, err
}

func (sh *SqlHandler) UpdateRecipe(id int, title string, makingTime string, serves string, ingredients string, cost int) ([]*constants.PatchRecipe, error) {
	// Prepare statement for inserting data
	stmt, err := sh.db.Prepare("UPDATE recipes set title = ?, making_time = ?, serves = ?, ingredients = ?, " +
		"cost = ? where id = ?") // ? = placeholder
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(title, makingTime, serves, ingredients, cost, id)
	if err != nil {
		return nil, err
	}
	if num, err := result.RowsAffected(); err != nil || num == 0 {
		return nil, errors.New(constants.RecipeNotFound)
	}
	recipe, err := sh.GetRecipeByID(id)
	var patchRecipe []*constants.PatchRecipe
	for _, v := range recipe {
		patchRecipe = append(patchRecipe, &constants.PatchRecipe{
			Title:       v.Title,
			MakingTime:  v.MakingTime,
			Serves:      v.Serves,
			Ingredients: v.Ingredients,
			Cost:        v.Cost,
		})
	}
	return patchRecipe, err
}

func (sh *SqlHandler) DeleteRecipeByID(id int) error {
	stmt, err := sh.db.Prepare("DELETE FROM recipes where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	if num, err := result.RowsAffected(); err != nil || num == 0 {
		return errors.New(constants.RecipeNotFound)
	}

	return nil
}
