package persistent

import (
	"Darkyfun/UrlShortener/internal/logging"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

var ErrConnect = errors.New("unable to connect")
var ErrNoRows = errors.New("no rows")
var ErrAlreadyExists = errors.New("alias already exists")
var ErrConnClosed = errors.New("connection closed")

type Db struct {
	pool Pooler
	log  logging.Logger
}

type Pooler interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
	Ping(ctx context.Context) error
}

func (d *Db) Close() {
	d.pool.Close()
}

func NewDb(ctx context.Context, logger *logging.EventLogger, conn string) *Db {
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		log.Fatalf("%s", ErrConnect)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	table := `create table if not exists url (
    alias varchar primary key ,
    original varchar,
    created_date timestamp
	);`

	_, _ = pool.Exec(ctx, table)

	return &Db{
		pool: pool,
		log:  logger,
	}
}

func (d *Db) GetOriginal(ctx context.Context, alias string) (string, error) {
	res := d.pool.QueryRow(ctx, `select original from url where alias = $1`, alias)
	var orig string

	err := res.Scan(&orig)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoRows
	}
	if err != nil && err.Error() == `closed pool` {
		d.log.Log("error", "unable to insert "+alias+" "+orig+" in sql: pool is closed")
		return "", ErrConnClosed
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "", ErrConnect
	}

	return orig, err
}

func (d *Db) GetAlias(ctx context.Context, orig string) (string, error) {
	res := d.pool.QueryRow(ctx, `select alias from url where original = $1`, orig)
	var alias string

	err := res.Scan(&alias)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoRows
	}
	if err != nil && err.Error() == `closed pool` {
		d.log.Log("error", "unable to insert "+alias+" "+orig+" in sql: pool is closed")
		return "", ErrConnClosed
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "", ErrConnect
	}

	return alias, nil
}

func (d *Db) Set(ctx context.Context, alias string, orig string) error {
	_, err := d.pool.Exec(ctx, `insert into url (alias, original, created_date) values ($1, $2, $3)`, alias, orig, time.Now())

	if err != nil && err.Error() == `ERROR: duplicate key value violates unique constraint "url_pkey" (SQLSTATE 23505)` {
		return ErrAlreadyExists
	}
	if err != nil && err.Error() == `closed pool` {
		d.log.Log("error", "unable to insert "+alias+" "+orig+" in sql: pool is closed")
		return ErrConnClosed
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrConnect
	}

	return err
}

func (d *Db) Ping(ctx context.Context) error {
	return d.pool.Ping(ctx)
}
