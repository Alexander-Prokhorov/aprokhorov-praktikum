package storage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"sync"
	"time"

	// used for pgx
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Pgs struct {
	DB         *sql.DB
	mutex      *sync.RWMutex
	buffer     Metrics
	bufferSize int
}

func NewDatabaseConnect(ctx context.Context, dbPath string) (Pgs, error) {
	const (
		defaultMaxIdleConns    = 10
		defaultMaxOpenConns    = 10
		defaultConnMaxIdleTime = 10
		defaultBufferSize      = 10
	)

	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return Pgs{}, err
	}

	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetConnMaxIdleTime(defaultConnMaxIdleTime)

	pgs := Pgs{
		DB:    db,
		mutex: &sync.RWMutex{},
		buffer: Metrics{
			Counter: make(map[string]Counter),
			Gauge:   make(map[string]Gauge),
		},
		bufferSize: defaultBufferSize,
	}

	err = pgs.InitDB(ctx)
	if err != nil {
		return Pgs{}, err
	}

	return pgs, nil
}

func (pgs *Pgs) InitDB(ctx context.Context) error {
	_, err := pgs.DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS Counter (name text PRIMARY KEY, value int8)")
	if err != nil {
		return err
	}

	_, err = pgs.DB.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS Gauge (name text  PRIMARY KEY, value float8)")
	if err != nil {
		return err
	}

	return nil
}

func (pgs Pgs) Ping(parentCtx context.Context) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*1)
	defer cancel()

	return pgs.DB.PingContext(ctx)
}

func (pgs Pgs) Read(ctx context.Context, valueType string, metricName string) (interface{}, error) {
	switch valueType {
	case "counter":
		return pgs.safeCounterRead(ctx, metricName)
	case "gauge":
		return pgs.safeGaugeRead(ctx, metricName)
	default:
		return nil, errors.New("PGS: Get(): Only [gauge, counter] type are supported")
	}
}

func (pgs Pgs) ReadAll(ctx context.Context) (map[string]map[string]string, error) {
	ret := make(map[string]map[string]string)
	ret["counter"] = make(map[string]string)
	ret["gauge"] = make(map[string]string)

	const (
		base    = 10
		bitSize = 64
	)

	var (
		name  string
		delta Counter
		value Gauge
	)

	// Get Values for Counter
	pgs.mutex.RLock()
	rows, err := pgs.DB.QueryContext(ctx, "SELECT * from Counter")
	pgs.mutex.RUnlock()

	if err != nil {
		return ret, err
	}

	if rows.Err() != nil {
		return ret, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&name, &delta)
		if err != nil {
			return ret, err
		}

		ret["counter"][name] = strconv.FormatInt(int64(delta), base)
	}

	// Get Values for Gauge
	pgs.mutex.RLock()
	rows, err = pgs.DB.QueryContext(ctx, "SELECT * from Gauge")
	pgs.mutex.RUnlock()

	if err != nil {
		return ret, err
	}

	if rows.Err() != nil {
		return ret, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err != nil {
			return ret, err
		}

		ret["gauge"][name] = strconv.FormatFloat(float64(value), 'f', -1, bitSize)
	}

	return ret, nil
}

func (pgs Pgs) Write(ctx context.Context, metricName string, value interface{}) error {
	switch data := value.(type) {
	case Counter:
		pgs.buffer.Counter[metricName] = data
	case Gauge:
		pgs.buffer.Gauge[metricName] = data
	default:
		err := errors.New("PGS: Post(): Only [gauge, counter] type are supported")

		return err
	}

	if len(pgs.buffer.Counter)+len(pgs.buffer.Gauge) == pgs.bufferSize {
		if err := pgs.Flush(ctx); err != nil {
			return err
		}
	}

	return pgs.Flush(ctx)
}

func (pgs *Pgs) safeCounterRead(ctx context.Context, metricName string) (Counter, error) {
	var value struct {
		Value int
		Valid bool
	}

	pgs.mutex.RLock()
	err := pgs.DB.QueryRowContext(ctx, "SELECT value from Counter where name = $1", metricName).Scan(&value.Value)
	pgs.mutex.RUnlock()

	if err != nil {
		return Counter(0), err
	}

	if value.Valid {
		return Counter(0), errors.New("value not found")
	}

	return Counter(value.Value), nil
}

func (pgs *Pgs) safeGaugeRead(ctx context.Context, metricName string) (Gauge, error) {
	var value struct {
		Value float64
		Valid bool
	}

	pgs.mutex.RLock()
	err := pgs.DB.QueryRowContext(ctx, "SELECT value from Gauge where name = $1", metricName).Scan(&value.Value)
	pgs.mutex.RUnlock()

	if err != nil {
		return Gauge(0), err
	}

	if value.Valid {
		return Gauge(0), errors.New("value not found")
	}

	return Gauge(value.Value), nil
}

func (pgs Pgs) Close() {
	pgs.DB.Close()
}

func (pgs Pgs) Flush(ctx context.Context) error {
	pgs.mutex.Lock()
	defer pgs.mutex.Unlock()

	tx, err := pgs.DB.Begin()
	if err != nil {
		return err
	}

	Stmt, err := pgs.DB.PrepareContext(
		ctx,
		"INSERT INTO Counter (name, value) "+
			"VALUES ($1, $2) "+
			"ON CONFLICT (name) "+
			"DO UPDATE SET value = EXCLUDED.value",
	)
	if err != nil {
		return err
	}

	txStmt := tx.StmtContext(ctx, Stmt)

	for metricName, metricValue := range pgs.buffer.Counter {
		_, err = txStmt.ExecContext(
			ctx,
			metricName,
			metricValue,
		)
		if err != nil {
			return err
		}
	}

	Stmt, err = pgs.DB.PrepareContext(
		ctx,
		"INSERT INTO Gauge (name, value) "+
			"VALUES ($1, $2) "+
			"ON CONFLICT (name) "+
			"DO UPDATE SET value = EXCLUDED.value",
	)
	if err != nil {
		return err
	}

	txStmt = tx.StmtContext(ctx, Stmt)

	for metricName, metricValue := range pgs.buffer.Gauge {
		_, err = txStmt.ExecContext(
			ctx,
			metricName,
			metricValue,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return err1
		}

		return err
	}

	return err
}
