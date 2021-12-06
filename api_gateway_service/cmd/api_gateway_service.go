package main

import (
	"api_gateway_service/pkg/api"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	// Структура сервера API Gateway
	type server struct {
		api *api.API
	}

	var srv server
	srv.api = api.New()

	// Запускаем сервис на порту 8080, интерфейс localhost.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("localhost:8080", srv.api.Router()))
	}()
	log.Println("API Gateway HTTP server started @ localhost:8080")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("API Gateway HTTP server stopped")
}
