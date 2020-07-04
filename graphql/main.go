package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/mitchellh/mapstructure"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

type GraphQLPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

var authors []Author = []Author{
	Author{
		ID:        "1",
		Firstname: "Igor",
		Lastname:  "Rashin",
		Username:  "Nish",
		Password:  "123",
	},
	Author{
		ID:        "2",
		Firstname: "Maria",
		Lastname:  "Rashin",
		Username:  "Ginger",
		Password:  "1234",
	},
}

var articles []Article = []Article{
	Article{
		ID:      "1",
		Author:  "1",
		Title:   "Article 1",
		Content: "Some awesome news",
	},
}

var rootQuery *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"authors": &graphql.Field{
			Type: graphql.NewList(authorType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return authors, nil
			},
		},
		"author": &graphql.Field{
			Type: authorType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				for _, author := range authors {
					if author.ID == id {
						return author, nil
					}
				}
				return nil, nil
			},
		},
		"articles": &graphql.Field{
			Type: graphql.NewList(articleType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return articles, nil
			},
		},
		"article": &graphql.Field{
			Type: articleType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				for _, article := range articles {
					if article.ID == id {
						return article, nil
					}
				}
				return nil, nil
			},
		},
	},
})

var rootMutation *graphql.Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"deleteAuthor": &graphql.Field{
			Type: graphql.NewList(authorType),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id := params.Args["id"].(string)
				for index, author := range authors {
					if author.ID == id {
						authors = append(authors[:index], authors[index+1:]...)
						return authors, nil
					}
				}
				return nil, nil
			},
		},
		"updateAuthor": &graphql.Field{
			Type: graphql.NewList(authorType),
			Args: graphql.FieldConfigArgument{
				"author": &graphql.ArgumentConfig{
					Type: authorInputType,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var data Author
				mapstructure.Decode(params.Args["author"], &data)
				for index, author := range authors {
					if author.ID == data.ID {
						if data.Firstname != "" {
							author.Firstname = data.Firstname
						}
						if data.Lastname != "" {
							author.Lastname = data.Lastname
						}
						if data.Username != "" {
							author.Username = data.Username
						}
						if data.Password != "" {
							hash, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
							author.Password = string(hash)
						}
						authors[index] = author
						return author, nil
					}
				}
				return nil, nil
			},
		},
		"createArticle": &graphql.Field{
			Type: graphql.NewList(articleType),
			Args: graphql.FieldConfigArgument{
				"article": &graphql.ArgumentConfig{
					Type: articleInputType,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var data Article
				mapstructure.Decode(params.Args["article"], &data)
				data.ID = uuid.Must(uuid.NewV4()).String()
				data.Author = "nish"
				articles = append(articles, data)
				return nil, nil
			},
		},
	}})

func main() {
	fmt.Println("app go")

	router := mux.NewRouter()
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	router.HandleFunc("/signin", Signin).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/graphql", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-type", "application/json")
		var payload GraphQLPayload
		json.NewDecoder(request.Body).Decode(&payload)
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  payload.Query,
			VariableValues: payload.Variables,
		})
		json.NewEncoder(response).Encode(result)
	})
	http.ListenAndServe(":12345", router)
}
