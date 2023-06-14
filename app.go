package main

import (
	"database/sql"

	// tom: for Initialize
	"fmt"
	"log"

	// tom: for route handlers
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/keploy/go-sdk/integrations/ksql/v2"

	// tom: go get required
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/miracle73/Go-Rest-API/model"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// tom: initial function is empty, it's filled afterwards
// func (a *App) Initialize(user, password, dbname string) { }
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":3200", a.Router))
	fmt.Printf("router initialized and listening on 3200\n")
}

// tom: added "sslmode=disable" to connection string
func (a *App) Initialize(user, password, dbname string) error {

	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", "5438", user, password, dbname)
	// connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	driver := ksql.Driver{Driver: pq.Driver{}}

	sql.Register("keploy", &driver)

	var err error
	a.DB, err = sql.Open("keploy", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	// tom: this line is added after initializeRoutes is created later on
	a.initializeRoutes()
	return err
}

// tom: these are added later
func (a *App) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Post ID")
		return
	}
	p := model.Post{ID: id}
	if err := p.GetPost(r.Context(), a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Post not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	Posts, err := model.GetAllPosts(r.Context(), a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Posts)
}

func (a *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	var p model.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.CreatePost(r.Context(), a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Post ID")
		return
	}

	var p model.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.UpdatePost(r.Context(), a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Post ID")
		return
	}

	p := model.Post{ID: id}
	if err := p.DeletePost(r.Context(), a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/posts", a.GetAllPosts).Methods("GET")
	a.Router.HandleFunc("/api/post/{id}", a.GetPost).Methods("GET")
	a.Router.HandleFunc("/api/post/new", a.CreatePost).Methods("POST")
	a.Router.HandleFunc("/api/post/update", a.UpdatePost).Methods("PUT")
	a.Router.HandleFunc("/api/post/delete/{id}", a.DeletePost).Methods("DELETE")
}
