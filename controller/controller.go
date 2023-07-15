package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/meinantoyuriawan/sharing-vison-backend/helper"
	"github.com/meinantoyuriawan/sharing-vison-backend/models"
	"gorm.io/gorm"
)

func CreateNewArticle(w http.ResponseWriter, r *http.Request) {
	userInput := models.Posts{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// check condition
	err := articleValidation(userInput)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	userInput.Created_date = time.Now()
	userInput.Updated_date = time.Now()

	// insert into db
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// response := map[string]string{"message": "success"}
	response := map[string]string{}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func ShowArticle(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	limit := params["limit"]
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		response := map[string]string{"message": "Limit error"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	offset := params["offset"]
	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		response := map[string]string{"message": "Offset error"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// access the database using user input

	var arrPost []models.Posts

	if err := models.DB.Limit(intLimit).Offset(intOffset).Find(&arrPost).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	response := arrPost
	helper.ResponseJSON(w, http.StatusOK, response)
}

func ShowArticleById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"message": "Id error"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	post := models.Posts{}
	if err := models.DB.Where("Id = ?", intId).First(&post).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}
	response := post
	helper.ResponseJSON(w, http.StatusInternalServerError, response)
}

func EditArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"message": "Id error"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// get user input
	userInput := models.Posts{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	// check condition
	err = articleValidation(userInput)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// get desired id
	post := models.Posts{}
	if err := models.DB.Where("Id = ?", intId).First(&post).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	post.Title = userInput.Title
	post.Content = userInput.Content
	post.Category = userInput.Category
	post.Status = userInput.Status
	post.Updated_date = time.Now()

	// push to db
	if err := models.DB.Save(&post).Error; err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// response := map[string]string{"message": "Edit success"}
	response := map[string]string{}
	helper.ResponseJSON(w, http.StatusOK, response)

}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	intId, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"message": "Id error"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	post := models.Posts{}
	if err := models.DB.Where("Id = ?", intId).First(&post).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// delete db
	models.DB.Delete(&post, intId)
	// response := map[string]string{"message": "Delete success"}
	response := map[string]string{}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func articleValidation(userInput models.Posts) error {

	titleLen := len([]rune(userInput.Title))
	contentLen := len([]rune(userInput.Content))
	categoryLen := len([]rune(userInput.Category))

	// title required minimum 20 char
	if titleLen < 20 {
		err := errors.New("title length less than 20 character")
		return err
	}
	// content required minimum 200 char
	if contentLen < 200 {
		err := errors.New("content length less than 200 character")
		return err
	}
	// category required minimum 3 char
	if categoryLen < 3 {
		err := errors.New("category length less than 3 character")
		return err
	}
	// status required between "publish" "draft" "thrash"
	if userInput.Status != "publish" {
		if userInput.Status != "draft" {
			if userInput.Status != "thrash" {
				err := errors.New("status Error")
				return err
			}
		}
	}
	return nil
}
