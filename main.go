package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"

	"github.com/mitchellh/mapstructure"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var JwtSecret []byte = []byte("bububu")

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

type CustomJWTClaim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func ValidateJWT(t string) (interface{}, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})

	if err != nil {
		return nil, errors.New(`{"message": "` + err.Error() + `" }`)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var tokenData CustomJWTClaim
		mapstructure.Decode(claims, &tokenData)
		return tokenData, nil
	}
	return nil, errors.New(`{"message": "invalid token" }`)
}

func RootRoute(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Write([]byte(`{"message": "Hello World"}`))
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				decoded, err := ValidateJWT(bearerToken[1])
				if err != nil {
					res.Header().Add("content-type", "application/json")
					res.WriteHeader(500)
					res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
					return
				}
				context.Set(req, "token", decoded)
				next(res, req)
			}
		} else {
			res.Header().Add("content-type", "application/json")
			res.WriteHeader(500)
			res.Write([]byte(`{ "message": "token is required" }`))
		}
	})
}

func main() {
	fmt.Println("Starting app...")

	router := mux.NewRouter()
	router.HandleFunc("/", RootRoute).Methods("GET")
	router.HandleFunc("/register", RegisterAuthor).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/authors", GetAllAuthors).Methods("GET")
	router.HandleFunc("/author/{id}", GetAuthor).Methods("GET")
	router.HandleFunc("/author/{id}", JWTMiddleware(DeleteAuthor)).Methods("DELETE")
	router.HandleFunc("/author/{id}", JWTMiddleware(UpdateAuthor)).Methods("PUT")
	router.HandleFunc("/articles", GetAllArticles).Methods("GET")
	router.HandleFunc("/article/{id}", GetArticle).Methods("GET")
	router.HandleFunc("/article/{id}", JWTMiddleware(DeleteArticle)).Methods("DELETE")
	router.HandleFunc("/article/{id}", JWTMiddleware(UpdateArticle)).Methods("PUT")
	router.HandleFunc("/article", JWTMiddleware(CreateArticle)).Methods("POST")
	methods := handlers.AllowedMethods(
		[]string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
	)
	headers := handlers.AllowedHeaders(
		[]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},
	)
	origins := handlers.AllowedOrigins(
		[]string{
			"*",
		},
	)

	http.ListenAndServe(":12345", handlers.CORS(headers, methods, origins)(router))
}
