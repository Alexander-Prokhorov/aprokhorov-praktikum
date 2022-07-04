package storage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Pgs struct {
	DB         *sql.DB
	mutex      *sync.RWMutex
	buffer     Metrics
	bufferSize int
	ctx        context.Context
}

func NewDatabaseConnect(dbPath string) (Pgs, error) {
	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return Pgs{}, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10)

	pgs := Pgs{
		DB:    db,
		mutex: &sync.RWMutex{},
		buffer: Metrics{
			Counter: make(map[string]Counter),
			Gauge:   make(map[string]Gauge),
		},
		bufferSize: 10,
		ctx:        context.Background(),
	}

	err = pgs.InitDB(pgs.ctx)
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
	if err := pgs.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (pgs Pgs) Read(valueType string, metricName string) (interface{}, error) {
	switch valueType {
	case "counter":
		return pgs.safeCounterRead(metricName)
	case "gauge":
		return pgs.safeGaugeRead(metricName)
	default:
		return nil, errors.New("PGS: Get(): Only [gauge, counter] type are supported")
	}
}

func (pgs Pgs) ReadAll() (map[string]map[string]string, error) {
	ret := make(map[string]map[string]string)
	ret["counter"] = make(map[string]string)
	ret["gauge"] = make(map[string]string)
	var (
		name  string
		delta Counter
		value Gauge
	)

	// Get Values for Counter
	pgs.mutex.RLock()
	rows, err := pgs.DB.Query("SELECT * from Counter")
	pgs.mutex.RUnlock()
	if err != nil {
		return ret, err
	}

	for rows.Next() {
		err = rows.Scan(&name, &delta)
		if err != nil {
			return ret, err
		}
		ret["counter"][name] = strconv.FormatInt(int64(delta), 10)
	}

	// Get Values for Gauge
	pgs.mutex.RLock()
	rows, err = pgs.DB.Query("SELECT * from Gauge")
	pgs.mutex.RUnlock()
	if err != nil {
		return ret, nil
	}

	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err != nil {
			return ret, err
		}
		ret["gauge"][name] = strconv.FormatFloat(float64(value), 'f', -1, 64)
	}

	return ret, nil
}

func (pgs Pgs) Write(metricName string, value interface{}) error {
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
		if err := pgs.Flush(); err != nil {
			return err
		}
	}

	if err := pgs.Flush(); err != nil {
		return err
	}

	return nil
}

func (pgs *Pgs) safeCounterRead(metricName string) (Counter, error) {
	var value struct {
		Value int
		Valid bool
	}

	pgs.mutex.RLock()
	err := pgs.DB.QueryRow("SELECT value from Counter where name = $1", metricName).Scan(&value.Value)
	pgs.mutex.RUnlock()
	if err != nil {
		return Counter(0), err
	}
	if value.Valid {
		return Counter(0), errors.New("value not found")
	}
	return Counter(value.Value), nil
}

func (pgs *Pgs) safeGaugeRead(metricName string) (Gauge, error) {
	var value struct {
		Value float64
		Valid bool
	}

	pgs.mutex.RLock()
	err := pgs.DB.QueryRow("SELECT value from Gauge where name = $1", metricName).Scan(&value.Value)
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

func (pgs Pgs) Flush() error {
	pgs.mutex.Lock()
	defer pgs.mutex.Unlock()

	tx, err := pgs.DB.Begin()
	if err != nil {
		return err
	}

	Stmt, err := pgs.DB.PrepareContext(pgs.ctx, "INSERT INTO Counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value")
	if err != nil {
		return err
	}

	txStmt := tx.StmtContext(pgs.ctx, Stmt)

	for metricName, metricValue := range pgs.buffer.Counter {
		_, err := txStmt.ExecContext(pgs.ctx,
			metricName,
			metricValue,
		)
		if err != nil {
			return err
		}
	}

	Stmt, err = pgs.DB.PrepareContext(pgs.ctx, "INSERT INTO Gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value")
	if err != nil {
		return err
	}

	txStmt = tx.StmtContext(pgs.ctx, Stmt)

	for metricName, metricValue := range pgs.buffer.Gauge {
		_, err := txStmt.ExecContext(pgs.ctx,
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
