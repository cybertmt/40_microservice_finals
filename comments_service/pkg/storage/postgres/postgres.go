package pgdb

import (
	"comments_service/pkg/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrorDuplicatePost error = errors.New("SQLSTATE 23505")

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

//New - Конструктор объекта хранилища.
func New(connstr string) (*Store, error) {

	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

//Close - освобождение БД
func (s *Store) Close() {
	s.db.Close()
}

//CommentsN получить все комментарии к новости n
func (s *Store) CommentsN(n int) ([]storage.Comment, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT 
		comments.id, 
		comments.author, 
		comments.content, 
		comments.pubtime, 
		comments.parentpost,
		comments.parentcomment
	FROM comments
	WHERE parentpost=$1;`, n)

	if err != nil {
		return nil, err
	}

	var comments []storage.Comment
	for rows.Next() {
		var c storage.Comment
		err = rows.Scan(
			&c.ID,
			&c.Author,
			&c.Content,
			&c.PubTime,
			&c.ParentPost,
			&c.ParentComment,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, rows.Err()
}

//AddComment - добавление комментария
func (s *Store) AddComment(c storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO comments (
		author, 
		content, 
		pubtime, 
		parentpost,
		parentcomment) 
	VALUES ($1,$2,$3,$4,$5);`,
		c.Author,
		c.Content,
		c.PubTime,
		c.ParentPost,
		c.ParentComment)
	return err
}

//UpdateComment - обновление комментария по id
func (s *Store) UpdateComment(c storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE comments 
	SET author=$2,
	content=$3,
	pubtime=$4,
	parentpost=$5,
	parentcomment=$6
	WHERE id=$1;`,
		c.ID,
		c.Author,
		c.Content,
		c.PubTime,
		c.ParentPost,
		c.ParentComment)
	return err
}

//DeleteComment - удаляет комментарий по id
func (s *Store) DeleteComment(c storage.Comment) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM comments 
	WHERE id=$1;`, c.ID)
	return err
}
