package persistent

import (
	"Darkyfun/UrlShortener/internal/logging"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"log"

	"io"
	"testing"
	"time"
)

const TestBase = "postgres://go:123@localhost:5432/service_test"

func TestMain(m *testing.M) {
	pool, err := pgxpool.New(context.Background(), TestBase)
	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.Exec(context.Background(), `drop table url`)
	if err != nil {
		log.Fatalf("failed to perform exec query 'drop table...' in test: %v\n", err)
	}

	_, err = pool.Exec(context.Background(), `create table if not exists url (
    	alias varchar primary key ,
    	original varchar,
    	created_date timestamp
		)`,
	)
	if err != nil {
		log.Fatalf("failed to perform exec query 'create table...' in test: %v\n", err)
	}

	_, err = pool.Exec(context.Background(), `insert into url (alias, original) values ('newone', 'testurl')`)
	if err != nil {
		log.Fatalf("failed to perform exec query 'insert into...' in test: %v\n", err)
	}

	// for redirect test.
	_, err = pool.Exec(context.Background(), `insert into url (alias, original) values ('googlealias', 'https://www.google.come')`)
	if err != nil {
		log.Fatalf("failed to perform exec query 'insert into...' in test: %v\n", err)
	}

	// before tests.
	m.Run()
	// after tests.
}

func TestDb_Set(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		origUrl string
		err     error
	}{
		{name: "new pair", alias: "new_alias", origUrl: "new_original_url", err: nil},
		{name: "already exists", alias: "new_alias", origUrl: "new_original_url", err: ErrAlreadyExists},
	}

	db := NewDb(context.Background(), logging.NewLogger("json", io.Discard), TestBase)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Set(ctx, tt.alias, tt.origUrl)
			assert.Equal(t, tt.err, err)
		})
	}

	ctxExp, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Nanosecond * 100)

	err := db.Set(ctxExp, "timeout", "timeout")
	assert.Equal(t, ErrConnect, err)

	// closing connection for catching ErrConnClosed.
	db.Close()

	err = db.Set(ctx, "closed connect", "closed connect")
	assert.Equal(t, ErrConnClosed, err)
}

func TestDb_GetOriginal(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		origUrl string
		err     error
	}{
		{name: "success", alias: "newone", origUrl: "testurl", err: nil},
		{name: "no rows", alias: "not_exist_alias", origUrl: "", err: ErrNoRows},
	}

	db := NewDb(context.Background(), logging.NewLogger("json", io.Discard), TestBase)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := db.GetOriginal(ctx, tt.alias)
			assert.Equal(t, err, tt.err)
			assert.Equal(t, res, tt.origUrl)
		})
	}

	ctxExp, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Nanosecond * 100)

	res, err := db.GetOriginal(ctxExp, "newone")
	assert.Equal(t, ErrConnect, err)
	assert.Equal(t, "", res)

	// closing connection for catching ErrConnClosed.
	db.Close()

	orig, err := db.GetOriginal(ctx, "closed connection")
	assert.Equal(t, "", orig)
	assert.Equal(t, ErrConnClosed, err)
}

func TestDb_GetAlias(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		origUrl string
		err     error
	}{
		{name: "success", alias: "newone", origUrl: "testurl", err: nil},
		{name: "no rows", alias: "", origUrl: "not_exist_url", err: ErrNoRows},
	}

	db := NewDb(context.Background(), logging.NewLogger("json", io.Discard), TestBase)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := db.GetAlias(ctx, tt.origUrl)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.alias, res)
		})
	}

	ctxExp, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Nanosecond * 100)

	res, err := db.GetAlias(ctxExp, "testurl")
	assert.Equal(t, ErrConnect, err)
	assert.Equal(t, "", res)

	db.Close()

	res, err = db.GetAlias(ctx, "closed connection")
	assert.Equal(t, ErrConnClosed, err)
	assert.Equal(t, "", res)
}
