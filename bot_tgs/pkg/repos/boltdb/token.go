package boltdb

import (
	"awesomeProject1/pkg/repos"
	"errors"
	"github.com/boltdb/bolt"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}

}

func (r *TokenRepository) Save(chatID int64, token string, bucket repos.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(inttobytes(chatID), []byte(token))
	})
}
func (r *TokenRepository) Get(chatID int64, bucket repos.Bucket) (string, error) {
	var token string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(inttobytes(chatID))
		token = string(data)
		return nil
	})
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", errors.New("Token not found.")
	}
	return token, nil
}

func inttobytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
