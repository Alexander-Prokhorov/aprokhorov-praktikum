package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Pgs struct {
	DB *sql.DB
}

func NewDatabaseConnect(dbPath string) (Pgs, error) {
	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return Pgs{}, err
	}
	return Pgs{DB: db}, nil
}

func (pgs Pgs) Ping(parentCtx context.Context) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*1)
	defer cancel()
	if err := pgs.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (pgs Pgs) Read(valueType string, metric string) (interface{}, error) {
	return nil, nil
}

func (pgs Pgs) ReadAll() map[string]map[string]string {
	return map[string]map[string]string{}
}

func (pgs Pgs) Write(metric string, value interface{}) error {
	return nil
}
