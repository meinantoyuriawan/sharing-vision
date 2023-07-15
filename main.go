package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meinantoyuriawan/sharing-vison-backend/controller"
	"github.com/meinantoyuriawan/sharing-vison-backend/models"
)

func main() {
	models.ConnectDB()
	r := mux.NewRouter()

	r.HandleFunc("/article", controller.CreateNewArticle).Methods("POST")
	r.HandleFunc("/article/{limit}/{offset}", controller.ShowArticle).Methods("GET")
	r.HandleFunc("/article/{id}", controller.ShowArticleById).Methods("GET")
	r.HandleFunc("/article/{id}", controller.EditArticle).Methods("POST")
	r.HandleFunc("/article/{id}", controller.DeleteArticle).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", r))
}
