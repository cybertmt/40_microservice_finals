package api

import (
	"comments_service/pkg/storage"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

const logfile = "log/comments_service_log.txt"

// Программный интерфейс сервера Comments
type API struct {
	db storage.Interface
	r  *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	// получить все комментарии к новости n
	api.r.HandleFunc("/comments/{n}", api.comments).Methods(http.MethodGet)
	// сохранить комментарий
	api.r.HandleFunc("/comments", api.storeComment).Methods(http.MethodPost)
	// сквозная идентификация запросов
	api.r.Use(api.reqId)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

// получить комментарии к новости n
func (api *API) comments(w http.ResponseWriter, r *http.Request) {
	ns := mux.Vars(r)["n"]
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := api.db.CommentsN(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// сохранение комментария
func (api *API) storeComment(w http.ResponseWriter, r *http.Request) {

	bComment, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("comsrv storeComment ReadAll error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c := storage.Comment{}
	err = json.Unmarshal(bComment, &c)
	if err != nil {
		http.Error(w, fmt.Sprintf("comsrv storeComment Unmarshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	c.PubTime = time.Now().Unix()

	err = api.db.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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