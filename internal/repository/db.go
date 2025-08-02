package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

// Client ...
type Client interface {
	Close() error
	DB() *DB // метод для подключения пулов к бд
}

// DB ...
type DB struct {
	pool *pgxpool.Pool
}

type Config struct { // конфиг, чтобы цепануться к слонику
	Login       string
	Password    string
	MasterHost  string
	ReplicaHost string
	Database    string
}

type client struct {
	db *DB
}

func NewClient(ctx context.Context, host string, cfg Config) (Client, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Login, cfg.Password, host, cfg.Database,
	)

	pool, err := pgxpool.New(ctx, url) // создаем пулы соединения к слонику (принимает контекст и url) (отдаем обьект пула и ошибку)
	if err != nil {
		if strings.Contains(err.Error(), cfg.Password) {
			err = errors.New(strings.ReplaceAll(err.Error(), cfg.Password, "xxxxxx"))
		}

		return nil, err
	}

	err = pool.Ping(ctx) // пингуем пул соединения
	if err != nil {
		return nil, err
	}

	return &client{
		db: &DB{pool: pool},
	}, nil
}

func (db *DB) Close() { // закрыть пул
	db.pool.Close()
}

func (c *client) Close() error {
	if c != nil {
		if c.db != nil {
			c.db.Close()
		}
	}

	return nil
}

func (c *client) DB() *DB { // возвращаем указатель на DB (вроде, чтобы работать с реальным пулом из структуры DB)
	return c.db
}
