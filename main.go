package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// apiHomePage defines the landing page for the whole API
func apiHomePage(w http.ResponseWriter, r *http.Request) {
	dataHomePage := "Endpoint: homepage"
	io.WriteString(w, dataHomePage)
}

// custom404PageHandler defines custom 404 page
func custom404PageHandler(w http.ResponseWriter, r *http.Request) {

	// set the content header type
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound) // this automatically generates a 404 status code
	// page content
	data404Page := "This page does not exist ... 404!"
	io.WriteString(w, data404Page)
}

// Article defines the data structure to be served via the API
type Article struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Content     string `json:"Content"`
}

// ArticlesStoreType is a new array declared as having elements which are of type Article
type ArticlesStoreType []Article

// mocking up some dummy data (as if you are extracting them from the pseudoDB Articles, into the variable articles)
var articles = ArticlesStoreType{
	Article{ID: "1", Title: "Love is blind", Description: "Here we discuss the generally held believe that love is blind", Content: "Article Content"},
	Article{ID: "2", Title: "Smooth Operator", Description: "This was the hit soulful music created by Sade Adu", Content: "Article Content"},
	Article{ID: "3", Title: "Jesus is Lord", Description: "Salvation is highly important - hence the need to see why Jesus is Lord", Content: "Article Content"},
	Article{ID: "4", Title: "Can money buy love", Description: "Here we explore whether Ladies indeed try to get married for money as the motivating factor", Content: "Article Content"},
}

// returnAllArticles() defines the request handling to return all artciles in the database
func returnAllArticles(w http.ResponseWriter, r *http.Request) {

	fmt.Println("JSON data is now being encoded ... ")
	json.NewEncoder(w).Encode(articles) // handles the json encoding part of the endpoint to same Writer used for html
}

// returnSingleArticle defines the request handling to return onlu articles specific to just one particular variable
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // extract the response object
	key := vars["id"]   // this gets the actual path variable that was entered

	// use the `key` to identify which json data to return by looping over the pesudoDB `Articles`
	for _, article := range articles { // iterate over the collection of `Article` structs
		if article.ID == key { // only encode for the object which the `ID` attribute is what was entered as the page variable
			json.NewEncoder(w).Encode(article)
		}
	}
	// json.NewEncoder(w).Encode(articles) // handles the json encoding part of the endpoint to same Writer used for html
}

// createnewArticle() defines the request handling creation of a new article
func creatNewArticle(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body) // extract the post content within the request body
	var article Article                      // initialize variable article meant to hold request body data
	json.Unmarshal(requestBody, &article)    // assumes that the request body JSON data, here we convert to Go bytes
	articles = append(articles, article)     // add the transformed data now to the articles variables
}

// deleteArticle() defines the request handling deletion of an existing article
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // extract the response object
	key := vars["id"]   // this gets the actual article identifier that was entered

	// use the `key` to identify which json data to delete by looping over the pesudoDB `Articles`
	for index, article := range articles { // iterate over the collection of `Article` structs
		if article.ID == key { // only encode for the object which the `ID` attribute is what was entered as the page variable
			articles = append(articles[:index], articles[index+1:]...) // note that ... is used to unpack elements of `articles[index+1]` before appending : because arguments are expected not another slice
		}
	}
}

// testPostArticles handles the editing and addition of articles to the API backend
func testPostArticle(w http.ResponseWriter, r *http.Request) {
	dataPOSTPage := "Endpoint: postpage"
	io.WriteString(w, dataPOSTPage)
}

// serviceRequestHandler() defines request handling service [used to aggregate all endpoints before running]
func serviceRequestHandlers() {
	muxRouter := mux.NewRouter().StrictSlash(true)                         // instantiate the gorillamux Router and enforce trailing slash rule i.e. `/path` === `/path/`
	muxRouter.NotFoundHandler = http.HandlerFunc(custom404PageHandler)     // customer 404 Page handler scenario
	muxRouter.HandleFunc("/", apiHomePage)                                 // declaring path and the handler function
	muxRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")    // responds only to GET methods i.e. selection [though it is same url]
	muxRouter.HandleFunc("/articles", testPostArticle).Methods("POST")     // responds only to POST methods i.e. insertion [though it is same url]
	muxRouter.HandleFunc("/article", creatNewArticle).Methods("POST")      // responds to POST request to create a new article
	muxRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE") // responds to DELETE request to remove an article with a specific identifier
	muxRouter.HandleFunc("/article/{id}", returnSingleArticle)             // responds to granular selection of articles
	log.Fatal(http.ListenAndServe(":9100", muxRouter))                     // set the port where the http server listens and serves. changed `nil` to the instance muxRouter
}

func main() {

	go serviceRequestHandlers() // call and run the server as a goroutine

	// create an artificial pause "so as to ensure the main functoin does not exit during goroutine runtime"
	var tempString string
	fmt.Scanln(&tempString)
}
