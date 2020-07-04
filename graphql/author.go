package main

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Author struct {
	ID        string `json:"id,omitempty" validate:"omitempty,uuid"`
	Firstname string `json:"first_name,omitempty" validate:"required"`
	Lastname  string `json:"last_name,omitempty" validate:"required"`
	Username  string `json:"user_name,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required,gte=4"`
}

var authorType *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"first_name": &graphql.Field{
			Type: graphql.String,
		},
		"last_name": &graphql.Field{
			Type: graphql.String,
		},
		"user_name": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var authorInputType *graphql.InputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "AuthorInput",
	Fields: &graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"first_name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"last_name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"user_name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"password": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

func Signin(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var author Author
	json.NewDecoder(request.Body).Decode(&author)
	hash, _ := bcrypt.GenerateFromPassword([]byte(author.Password), 10)
	author.ID = uuid.Must(uuid.NewV4()).String()
	author.Password = string(hash)
	authors = append(authors, author)
	json.NewEncoder(response).Encode(authors)
}

func Login(response http.ResponseWriter, req *http.Request) {
	response.Header().Add("content-type", "application/json")
	var data Author
	json.NewDecoder(req.Body).Decode(&data)
	for _, author := range authors {
		if author.Username == data.Username {
			err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(data.Password))
			if err != nil {
				response.WriteHeader(500)
				response.Write([]byte(`{"message": "invalid password"}`))
				return
			}
			json.NewEncoder(response).Encode(author)
		}
		return
	}
	json.NewEncoder(response).Encode(Author{})
}
