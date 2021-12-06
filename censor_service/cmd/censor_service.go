package main

import (
	"cens_service/pkg/api"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	// Структура сервера censor.
	type server struct {
		api *api.API
	}

	var srv server

	srv.api = api.New()

	// Запускаем сервис на порту 8083, интерфейс localhost.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("localhost:8083", srv.api.Router()))
	}()
	log.Println("Censor HTTP server started @ localhost:8083")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("Censor HTTP server stopped")
}
