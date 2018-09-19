package possum

import (
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

var (
	_pool                 *pgx.ConnPool
	ErrPoolNotInitialized = errors.New("Connection pool is not initialized")
)

type DB struct {
	ConnPool *pgx.ConnPool
}

func NewDB(host string, port int, user string, dbname string, password string, sslMode string, maxConnections int) *DB {

	db := &DB{}

	if _pool == nil {

		connConfig := pgx.ConnConfig{
			Host:      host,
			Port:      uint16(port),
			User:      user,
			Database:  dbname,
			Password:  password,
			LogLevel:  5,
			TLSConfig: nil,
		}

		poolConfig := pgx.ConnPoolConfig{
			ConnConfig:     connConfig,
			MaxConnections: maxConnections,
			AcquireTimeout: 0,
		}

		pool, err := pgx.NewConnPool(poolConfig)

		if err != nil {
			log.Fatal(err)
		}

		_pool = pool

	}

	db.ConnPool = _pool

	return db

}

func GetConnPool() (*pgx.ConnPool, error) {

	if _pool != nil {
		return _pool, nil
	}

	return nil, ErrPoolNotInitialized

}

func GetConnectionForSchema(schema string) (*pgx.ConnPool, *pgx.Conn) {

	pool, err := GetConnPool()

	if err != nil {
		panic(err)
	}

	conn, err := pool.Acquire()

	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(fmt.Sprintf("set search_path to %s", schema))

	if err != nil {
		panic(err)
	}

	return pool, conn

}

func Stat() pgx.ConnPoolStat {

	pool, err := GetConnPool()

	if err != nil {
		panic(err)
	}

	return pool.Stat()

}

func Query(schema string, query string, args pgx.QueryArgs) (*pgx.Rows, error) {

	pool, conn := GetConnectionForSchema(schema)

	defer pool.Release(conn)

	return conn.Query(query, args...)

}

func QueryRow(schema, query string, args pgx.QueryArgs) *pgx.Row {

	pool, conn := GetConnectionForSchema(schema)

	defer pool.Release(conn)

	return conn.QueryRow(query, args...)

}

func Create(schema string, out interface{}) (string, error) {

	var uuid pgtype.UUID
	scope := NewScope(out)

	err := QueryRow(schema, scope.CreateSQL(), scope.CreateArgs()).Scan(&uuid)

	if err != nil {
		return "", err
	}

	return EncodeUUID(uuid.Bytes), nil

}
