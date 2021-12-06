package pgdb

import (
	"comments_service/pkg/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		connstr string
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.connstr)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_AddComment(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		c storage.Comment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			if err := s.AddComment(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("AddComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_Close(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestStore_CommentsN(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []storage.Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			got, err := s.CommentsN(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentsN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentsN() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_DeleteComment(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		c storage.Comment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			if err := s.DeleteComment(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_UpdateComment(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		c storage.Comment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			if err := s.UpdateComment(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("UpdateComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
