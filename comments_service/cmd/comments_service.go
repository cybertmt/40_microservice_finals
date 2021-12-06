package main

import (
	"comments_service/pkg/api"
	"comments_service/pkg/storage"
	"comments_service/pkg/storage/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	// Структура сервера comments.
	type server struct {
		db  storage.Interface
		api *api.API
	}

	// Создаём объект сервера
	var srv server

	//  Создаём объект базы данных PostgreSQL.
	host := os.Getenv("host")
	bdName2 := os.Getenv("bdName2") // название БД для комментариев
	pwd := os.Getenv("pwd")
	db, err := pgdb.New("postgres://cyber:" + pwd + "@" + host + "/" + bdName2)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера БД
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем сервис на порту 8082, интерфейс localhost.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("localhost:8082", srv.api.Router()))
	}()
	log.Println("Comments HTTP server started @ localhost:8082")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("Comments HTTP server stopped")
}
