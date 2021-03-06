package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"news_service/pkg/storage"
	"os"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

const logfile = "log/news_service_log.txt"

// Программный интерфейс сервера News
type API struct {
	db     storage.Interface
	router *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	// получить страницу [page] последних новостей с указанием количества [quantity] на страницу
	api.router.HandleFunc("/news/{page}/{quantity}", api.posts).Methods(http.MethodGet)
	// получить страницу [page] последних новостей с указанием количества [quantity] на страницу
	// с [keyword] в заголовке новости
	api.router.HandleFunc("/filter/{page}/{quantity}/{keyword}", api.filter).Methods(http.MethodGet)
	// получить детальную новость по n
	api.router.HandleFunc("/detailed/{n}", api.detailed).Methods(http.MethodGet)
	// сквозная идентификация запросов
	api.router.Use(api.reqId)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

//  получить страницу [page] последних новостей с указанием количества [quantity] на страницу
func (api *API) posts(w http.ResponseWriter, r *http.Request) {

	ns := mux.Vars(r)["page"]
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	qs := mux.Vars(r)["quantity"]
	q, err := strconv.Atoi(qs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type result struct {
		Count int
		Posts []storage.Post
	}

	var res result
	res.Posts, res.Count, err = api.db.PostsN(n, q)
	if err != nil {
		http.Error(w, fmt.Sprintf("PostsN error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// получить страницу [page] последних новостей с указанием количества [quantity] на страницу
// с [keyword] в заголовке новости
func (api *API) filter(w http.ResponseWriter, r *http.Request) {

	ns := mux.Vars(r)["page"]
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	qs := mux.Vars(r)["quantity"]
	q, err := strconv.Atoi(qs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := mux.Vars(r)["keyword"]

	type result struct {
		Count int
		Posts []storage.Post
	}

	var res result
	res.Posts, res.Count, err = api.db.Filter(n, q, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Filter error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// получить детальную новость по n
func (api *API) detailed(w http.ResponseWriter, r *http.Request) {

	ns := mux.Vars(r)["n"]
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p storage.Post
	p, err = api.db.Post(n)
	if err != nil {
		http.Error(w, fmt.Sprintf("Post error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// сквозная идентификация запросов
// ?request_id=112233445566
func (api *API) reqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			http.Error(w, fmt.Sprintf("os.OpenFile error: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		id := r.URL.Query().Get("request_id")
		if id == "" {
			uid, err := uuid.NewV4()
			if err != nil {
				http.Error(w, fmt.Sprintf("uuid.NewV4 error: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			id = uid.String()
		}
		ctx := context.WithValue(r.Context(), "request_id", id)
		r = r.WithContext(ctx)
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		for k, v := range rec.Result().Header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)

		fmt.Fprintf(file, "Request ID: %s\n", id)
		fmt.Fprintf(file, "Time: %s\n", time.Now().Format(time.RFC1123))
		fmt.Fprintf(file, "Remote IP: %s\n", r.RemoteAddr)
		fmt.Fprintf(file, "HTTP Status: %d\n", rec.Result().StatusCode)
		fmt.Fprintln(file)
	})
}