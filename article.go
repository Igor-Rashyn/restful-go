package main

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/mux"
)

//Article model
type Article struct {
	ID      string `json:"id,omitempty"`
	Author  string `json:"author,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
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
	for index, article := range articles {
		if article.ID == params["id"] {
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
	var newArticle Article
	json.NewDecoder(req.Body).Decode(&newArticle)
	for index, article := range articles {
		if article.ID == params["id"] {
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
	article.ID = uuid.Must(uuid.NewV4()).String()
	articles = append(articles, article)
	json.NewEncoder(res).Encode(articles)
}
