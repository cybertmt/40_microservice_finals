package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

const logfile = "log/censor_service_log.txt"

// Программный интерфейс приложения
type API struct {
	r *mux.Router
}

// Конструктор объекта API
func New() *API {
	api := API{}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	// проверка комментария
	api.r.HandleFunc("/censor", api.censor).Methods(http.MethodPost)
	// сковзная идентификация запросов
	api.r.Use(api.reqId)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

// Список запрещенных слов.
var fWords = []string{"qwerty", "йцукен", "zxvbnm"}

// Проверка наличия подстрок [subs] в строке [str]
func censored(str string, subs ...string) bool {
	censored := false
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			censored = true
		}
	}
	return censored
}

// проверка комментария на наличие запрещенных слов
// localhost:8083/censor
func (api *API) censor(w http.ResponseWriter, r *http.Request) {
	bComment, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("censor ReadAll error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if censored(string(bComment), fWords...) {
		http.Error(w, "Censored", http.StatusNotAcceptable)
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