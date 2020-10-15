package main

import (
	"RestAPIForRecipe/constants"
	"RestAPIForRecipe/db"
	"RestAPIForRecipe/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var recipes = "recipes"
var sqlHandler = db.NewSqlHandler()

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		result, err := sqlHandler.GetRecipes()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error. DB request failed.\n")
			return
		}
		recipes := convertRecipeToRecipeWithoutTime(result)
		res, err := json.Marshal(&constants.GetRecipesResponse{Recipe: recipes})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(res))
		return
	case http.MethodPost:
		body := r.Body
		defer body.Close()
		param, err := getAndCheckParams(body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			s, err := json.Marshal(&constants.ErrorPostRecipesResponse{
				Message:  constants.ErrorCreationFailed,
				Required: err.Error(),
			})
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, "internal server error.\n")
				return
			}
			fmt.Fprintf(w, string(s))
			return
		}
		recipe, err := sqlHandler.InsertRecipe(param.Title, param.MakingTime, param.Serves, param.Ingredients, param.Cost)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error. insert failed.\n")
			return
		}

		s, err := json.Marshal(constants.PostRecipesResponse{
			Message: constants.RecipeCreateSuccess,
			Recipe:  recipe,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error.\n")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(s))
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "request method not allowed.\n")
	}
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	tmp := r.URL.Path[len("/"+recipes+"/"):]
	id, err := strconv.Atoi(tmp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request. %s is not valid ID.\n", tmp)
		return
	}

	switch r.Method {
	case http.MethodGet:
		result, err := sqlHandler.GetRecipeByID(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error. DB request failed.\n")
			return
		}
		recipes := convertRecipeToRecipeWithoutTime(result)
		res, err := json.Marshal(&constants.GetRecipeByIDResponse{Message: constants.GetRecipeByIDMessage,
			Recipe: recipes})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(res))
		return
	case http.MethodPatch:
		body := r.Body
		defer body.Close()
		param, err := getAndCheckParams(body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			s, err := json.Marshal(&constants.ErrorPostRecipesResponse{
				Message:  constants.ErrorCreationFailed,
				Required: err.Error(),
			})
			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, "internal server error.\n")
				return
			}
			fmt.Fprintf(w, string(s))
			return
		}
		recipe, err := sqlHandler.UpdateRecipe(id, param.Title, param.MakingTime, param.Serves, param.Ingredients,
			param.Cost)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			s, _ := json.Marshal(&constants.MessageResponse{Message: err.Error()})
			fmt.Fprintf(w, string(s))
			return
		}

		s, err := json.Marshal(constants.PatchRecipesResponse{
			Message: constants.RecipeUpdateSuccess,
			Recipe:  recipe,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error.\n")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(s))
		return
	case http.MethodDelete:
		body := r.Body
		defer body.Close()
		err := sqlHandler.DeleteRecipeByID(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			s, _ := json.Marshal(&constants.MessageResponse{Message: err.Error()})
			fmt.Fprintf(w, string(s))
			return
		}

		s, err := json.Marshal(constants.MessageResponse{
			Message: constants.RecipeDeleteSuccess,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error.\n")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(s))
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "request method not allowed.\n")
	}

}

func convertRecipeToRecipeWithoutTime(recipes []*models.Recipe) []*constants.RecipeWithoutTime {
	var withoutTime []*constants.RecipeWithoutTime
	for _, v := range recipes {
		withoutTime = append(withoutTime, &constants.RecipeWithoutTime{
			ID:          v.ID,
			Title:       v.Title,
			MakingTime:  v.MakingTime,
			Serves:      v.Serves,
			Ingredients: v.Ingredients,
			Cost:        v.Cost,
		})
	}
	return withoutTime
}

func getAndCheckParams(body io.ReadCloser) (*constants.RecipeRequest, error) {
	var recipe constants.RecipeRequest
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	json.Unmarshal(buf.Bytes(), &recipe)
	if recipe.Title == "" || recipe.MakingTime == "" || recipe.Serves == "" || recipe.Ingredients == "" ||
		recipe.Cost == 0 {
		return nil, errors.New(constants.ErrorParameterRequired)
	}
	return &recipe, nil
}

func main() {
	http.HandleFunc("/"+recipes, handler)
	http.HandleFunc("/"+recipes+"/", idHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
