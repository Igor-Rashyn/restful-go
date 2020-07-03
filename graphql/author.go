package main

import (
	"github.com/graphql-go/graphql"
)

type Author struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"first_name,omitempty"`
	Lastname  string `json:"last_name,omitempty"`
	Username  string `json:"user_name,omitempty"`
	Password  string `json:"password,omitempty"`
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
