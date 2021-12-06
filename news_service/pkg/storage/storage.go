package storage

// Post - публикация.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubDate string // дата публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Post(int) (Post, error)                       // получение одной новости
	PostsN(int, int) ([]Post, int, error)         // получение страницы новостей
	Filter(int, int, string) ([]Post, int, error) // получение отфильтрованной страницы новостей
	AddPost(Post) error                           // добавление новости
	UpdatePost(Post) error                        // обновление новости
	DeletePost(Post) error                        // удаление новости по ID
	Close()                                       // закрытие БД
}

