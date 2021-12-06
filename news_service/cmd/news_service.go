package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"news_service/pkg/api"
	"news_service/pkg/rss"
	"news_service/pkg/storage"
	"news_service/pkg/storage/postgres"
	"os"
	"os/signal"
)

func main() {

	// канал для чтения новостей rss-каналов и записи в бд
	PostChannel := make(chan storage.Post)
	// канал сбора ошибок rss обработчика
	ErrorChannel := make(chan error)

	// структура config.jason для rss обработчика
	type rssConfig struct {
		Rss     []string `json:"rss"`
		RPeriod int64      `json:"request_period"`
	}

	// Читаем файл конфигурации в последовательность байт
	jsonFile, err := os.Open("cmd/config.json")
	if err != nil {
		ErrorChannel <- err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Конвертируем байты в структуру rss конфигурации
	var rssCfg rssConfig
	_ = json.Unmarshal(byteValue, &rssCfg)

	// Запускаем rss обработчики в отдельных горутинах с конфигурациями
	for i := 0; i < len(rssCfg.Rss); i++ {
		go rss.FetchRss(rssCfg.Rss[i], rssCfg.RPeriod, PostChannel, ErrorChannel)
	}

	// Структура сервера news.
	type server struct {
		db  storage.Interface
		api *api.API
	}
	// Создаём объект сервера.
	var srv server

	// используем переменные окружения для адреса бд, имени бд и пароля
	host := os.Getenv("host")
	bdName1 := os.Getenv("bdName1") // название БД для новостей
	pwd := os.Getenv("pwd")

	// Реляционная БД PostgreSQL.
	db, err := postgres.New("postgres://cyber:" + pwd + "@" + host + "/" + bdName1)
	if err != nil {
		ErrorChannel <- err
	}

	//Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	defer srv.db.Close()

	// Запускаем запись новостей из канала PostChannel в БД в отдельной горутине
	go func() {
		for {
			err = srv.db.AddPost(<- PostChannel)
			if err != nil {
				ErrorChannel <- err
			}
		}
	}()

	// Запускаем логирование ошибок из канала ErrorChannel
	go func() {
		for {
			log.Println(<-ErrorChannel)
		}
	}()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем сервис на порту 8081, интерфейс localhost.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("localhost:8081", srv.api.Router()))
	}()
	log.Println("News HTTP server started @ localhost:8081")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("News HTTP server stopped")
}
