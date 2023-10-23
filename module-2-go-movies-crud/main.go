package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:id`
	Isbn     string    `json:isbn`
	Title    string    `json:title`
	Director *Director `json:director`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// create a slice of movies
var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

// Most difficult function: because we remove the actual movie and readd the movie (only in this no database example)
func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set json content type
	w.Header().Set("Content-Type", "application/json")
	// params
	params := mux.Vars(r)
	// loop over the movies, range
	// delete the movie with the id what you've sent
	// add a new movie - the movie that we senen in the body of postman
	for index, item := range movies {
		if item.ID == params["id"] {
			// delete in the append
			movies = append(movies[:index], movies[index+1:]...)
			// prepare a new movie variable
			var movie Movie
			// decode from the body into an struct directly
			_ = json.NewDecoder(r.Body).Decode(&movie)
			// set the ID
			movie.ID = params["id"]
			// add the movie into the movies appending
			movies = append(movies, movie)
			// write movie
			json.NewEncoder(w).Encode(movie)
			return
		}

	}
}

func main() {
	r := mux.NewRouter()

	// Append two particular movies and directors to the "movies" slice
	movies = append(movies, Movie{ID: "1", Isbn: "334564", Title: "Movie Fatal", Director: &Director{Firstname: "John", Lastname: "Salchichon"}})
	movies = append(movies, Movie{ID: "2", Isbn: "787746", Title: "Movie Alive", Director: &Director{Firstname: "Alpha", Lastname: "Golf"}})
	r.HandlerFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
