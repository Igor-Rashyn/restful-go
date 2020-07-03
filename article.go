package main

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Article model
type Article struct {
	ID      string `json:"id,omitempty" validate:"omitempty,uuid"`
	Author  string `json:"author,omitempty" validate:"isdefault"`
	Title   string `json:"title,omitempty" validate:"required"`
	Content string `json:"content,omitempty" validate:"required"`
}

func GetAllArticles(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	json.NewEncoder(res).Encode(articles)
}

func GetArticle(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	for _, article := range articles {
		if article.ID == params["id"] {
			json.NewEncoder(res).Encode(article)
			return
		}
	}
	json.NewEncoder(res).Encode(Article{})
}

func DeleteArticle(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	token := context.Get(req, "token").(CustomJWTClaim)
	for index, article := range articles {
		if article.ID == params["id"] && article.Author == token.ID {
			articles = append(articles[:index], articles[index+1:]...)
			json.NewEncoder(res).Encode(articles)
			return
		}
	}
	json.NewEncoder(res).Encode(Article{})
}

func UpdateArticle(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	token := context.Get(req, "token").(CustomJWTClaim)
	var newArticle Article
	json.NewDecoder(req.Body).Decode(&newArticle)
	for index, article := range articles {
		if article.ID == params["id"] && article.Author == token.ID {
			if newArticle.Title != "" {
				article.Title = newArticle.Title
			}
			if newArticle.Content != "" {
				article.Content = newArticle.Content
			}
			articles[index] = article
			json.NewEncoder(res).Encode(articles)
			return
		}
	}
	json.NewEncoder(res).Encode(Article{})
}

func CreateArticle(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var article Article
	json.NewDecoder(req.Body).Decode(&article)
	token := context.Get(req, "token").(CustomJWTClaim)
	validate := validator.New()
	err := validate.Struct(article)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	article.ID = uuid.Must(uuid.NewV4()).String()
	article.Author = token.ID
	articles = append(articles, article)
	json.NewEncoder(res).Encode(articles)
}
