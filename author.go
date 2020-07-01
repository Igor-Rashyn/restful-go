package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// Author model
type Author struct {
	ID        string `json:"id,omitempty" validate:"omitempty,uuid"`
	Firstname string `json:"first_name,omitempty" validate:"required"`
	Lastname  string `json:"last_name,omitempty" validate:"required"`
	Username  string `json:"user_name,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required,gte=4"`
}

func Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var data Author
	json.NewDecoder(req.Body).Decode(&data)
	validate := validator.New()
	err := validate.StructExcept(data, "Firstname", "Lastname")
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	for _, author := range authors {
		if author.Username == data.Username {
			err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(data.Password))
			if err != nil {
				res.WriteHeader(500)
				res.Write([]byte(`{"message": "invalid password"}`))
				return
			}
		}
		json.NewEncoder(res).Encode(author)
		return
	}
	res.Write([]byte(`{"message": "invalid user name"}`))
}

func RegisterAuthor(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	var newAuthor Author
	json.NewDecoder(req.Body).Decode(&newAuthor)
	validate := validator.New()
	err := validate.Struct(newAuthor)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newAuthor.Password), 10)
	newAuthor.ID = uuid.Must(uuid.NewV4()).String()
	newAuthor.Password = string(hash)
	authors = append(authors, newAuthor)
	json.NewEncoder(res).Encode(authors)
}

func GetAllAuthors(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	json.NewEncoder(res).Encode(authors)
}

func GetAuthor(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	for _, author := range authors {
		if author.ID == params["id"] {
			json.NewEncoder(res).Encode(author)
			return
		}
	}
	json.NewEncoder(res).Encode(Author{})
}

func DeleteAuthor(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	for index, author := range authors {
		if author.ID == params["id"] {
			authors = append(authors[:index], authors[index+1:]...)
			json.NewEncoder(res).Encode(authors)
			return
		}
	}
	json.NewEncoder(res).Encode(Author{})
}

func UpdateAuthor(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	var newAuthor Author
	json.NewDecoder(req.Body).Decode(&newAuthor)
	validate := validator.New()
	err := validate.StructExcept(newAuthor, "Firstname", "Lastname", "Username", "Password")
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	for index, author := range authors {
		if author.ID == params["id"] {
			if newAuthor.Firstname != "" {
				author.Firstname = newAuthor.Firstname
			}
			if newAuthor.Lastname != "" {
				author.Lastname = newAuthor.Lastname
			}
			if newAuthor.Username != "" {
				author.Username = newAuthor.Username
			}
			if newAuthor.Password != "" {
				err = validate.Var(newAuthor.Password, "gte=4")
				if err != nil {
					res.WriteHeader(500)
					res.Write([]byte(`{ "message": "` + err.Error() + `"}`))
					return
				}
				hash, _ := bcrypt.GenerateFromPassword([]byte(newAuthor.Password), 10)
				author.Password = string(hash)
			}
			authors[index] = author
			json.NewEncoder(res).Encode(authors)
			return
		}
	}
	json.NewEncoder(res).Encode(Author{})
}
