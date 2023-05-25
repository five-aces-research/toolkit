package pg_wrapper

import (
	"context"
	"fmt"
	"github.com/five-aces-research/toolkit/fas"
	"github.com/five-aces-research/toolkit/fas/bybit"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper/qq"
	_ "github.com/golang-migrate/migrate/v4/source/file" // this is needed to use file:// URLs
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"strings"
)

type Pgx struct {
	db *pgxpool.Pool
	q  *qq.Queries
}

var ctx = context.Background()

func Connect(host, dbName, user, password string, port int) (*Pgx, error) {
	params := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	conn, err := pgxpool.New(ctx, params)
	if err != nil {
		return nil, err
	}

	Q := qq.New(conn)
	return &Pgx{
		db: conn,
		q:  Q,
	}, nil
}

func (pg *Pgx) Ping() error {
	return pg.db.Ping(context.Background())
}

func (pg *Pgx) GetTickers() ([]qq.GetAvaibleTickersRow, error) {
	return pg.q.GetAvaibleTickers(ctx)
}

func getExchangeId(name string) int32 {
	switch strings.ToUpper(name) {
	case "BYBIT":
		return 1
	default:
		log.Panicln("not implemented")
		return 0
	}
}

func loadExchanger(id int32) fas.Public {
	switch id {
	case 1:
		return bybit.NewPublic()
	default:
		return nil
	}
}
