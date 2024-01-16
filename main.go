package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	prolog   *Prolog
	database SimpleRepository
)

func main() {
	fmt.Println("Hello, World!")
	initDependencies()
	defer database.Close()

	handler := http.NewServeMux()
	handler.HandleFunc("/users", usersHandler)
	handler.HandleFunc("/tags", tagsHandler)

	err := http.ListenAndServe(":8080", handler)

	if err != nil {
		panic(err)
	}
}

func initDependencies() {
	prolog = NewProlog()

	primary, err := NewSQLiteRepository()
	secondary := &SlowDatabase{SimpleRepository: primary}
	database = &CuratorDatabase{SimpleRepository: secondary, curator: prolog}
	// database = &SlowDatabase{SimpleRepository: primary}

	if err != nil {
		panic(err)
	}

	err = Migrate(primary)

	if err != nil {
		panic(err)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch getMethod(r) {
	case "GET":
		getUsersHandler(w, r)
	case "POST":
		postUsersHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsersWithTags()

	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(users)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

func postUsersHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &User{Name: name}

	err := database.SaveUser(user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(name))
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	switch getMethod(r) {
	case "POST":
		postTagsHandler(w, r)
	case "DELETE":
		deleteTagsHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getMethod(r *http.Request) string {
	method := r.Method

	if r.URL.Query().Get("method") != "" {
		method = strings.ToUpper(r.URL.Query().Get("method"))
	}

	return method
}

func postTagsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	name := r.URL.Query().Get("name")

	if userID == "" || name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tag := &Tag{UserID: user.ID, Name: name}

	err = database.SaveTag(tag)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(name))
}

func deleteTagsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(id)

	if user == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tags, err := database.GetUserTags(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, tag := range tags {
		err = database.DeleteTag(&tag)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(userId))
}
