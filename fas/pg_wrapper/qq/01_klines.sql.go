// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: 01_klines.sql

package qq

import (
	"context"
)

const createTicker = `-- name: CreateTicker :one
INSERT INTO tickers(exchange_id, name) VALUES ($1,$2) RETURNING ticker_id
`

type CreateTickerParams struct {
	ExchangeID int32  `json:"exchange_id"`
	Name       string `json:"name"`
}

func (q *Queries) CreateTicker(ctx context.Context, arg CreateTickerParams) (int32, error) {
	row := q.db.QueryRow(ctx, createTicker, arg.ExchangeID, arg.Name)
	var ticker_id int32
	err := row.Scan(&ticker_id)
	return ticker_id, err
}

const getTickerId = `-- name: GetTickerId :one
SELECT ticker_id FROM tickers WHERE exchange_id = $1 and name = $2
`

type GetTickerIdParams struct {
	ExchangeID int32  `json:"exchange_id"`
	Name       string `json:"name"`
}

func (q *Queries) GetTickerId(ctx context.Context, arg GetTickerIdParams) (int32, error) {
	row := q.db.QueryRow(ctx, getTickerId, arg.ExchangeID, arg.Name)
	var ticker_id int32
	err := row.Scan(&ticker_id)
	return ticker_id, err
}

const getTickers = `-- name: GetTickers :many
SELECT name FROM tickers
`

func (q *Queries) GetTickers(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getTickers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const minAndMax = `-- name: MinAndMax :one
SELECT COALESCE(MAX(starttime),0)::BIGINT as Max, COALESCE(Min(starttime),0)::BIGINT as Min FROM klines
WHERE ticker_id = $1 and resolution = $2
`

type MinAndMaxParams struct {
	TickerID   int32 `json:"ticker_id"`
	Resolution int32 `json:"resolution"`
}

type MinAndMaxRow struct {
	Max int64 `json:"max"`
	Min int64 `json:"min"`
}

func (q *Queries) MinAndMax(ctx context.Context, arg MinAndMaxParams) (MinAndMaxRow, error) {
	row := q.db.QueryRow(ctx, minAndMax, arg.TickerID, arg.Resolution)
	var i MinAndMaxRow
	err := row.Scan(&i.Max, &i.Min)
	return i, err
}

const readOHCL = `-- name: ReadOHCL :many
SELECT starttime, open, high, close, low, volume
FROM klines
WHERE ticker_id = $1
  AND resolution = $2
  AND starttime > $3
  AND starttime < $4
ORDER BY starttime asc
`

type ReadOHCLParams struct {
	TickerID   int32 `json:"ticker_id"`
	Resolution int32 `json:"resolution"`
	St         int64 `json:"st"`
	Et         int64 `json:"et"`
}

type ReadOHCLRow struct {
	Starttime int64   `json:"starttime"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Close     float64 `json:"close"`
	Low       float64 `json:"low"`
	Volume    float64 `json:"volume"`
}

func (q *Queries) ReadOHCL(ctx context.Context, arg ReadOHCLParams) ([]ReadOHCLRow, error) {
	rows, err := q.db.Query(ctx, readOHCL,
		arg.TickerID,
		arg.Resolution,
		arg.St,
		arg.Et,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadOHCLRow
	for rows.Next() {
		var i ReadOHCLRow
		if err := rows.Scan(
			&i.Starttime,
			&i.Open,
			&i.High,
			&i.Close,
			&i.Low,
			&i.Volume,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type WriteOHCLVParams struct {
	TickerID   int32   `json:"ticker_id"`
	Resolution int32   `json:"resolution"`
	Starttime  int64   `json:"starttime"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Close      float64 `json:"close"`
	Low        float64 `json:"low"`
	Volume     float64 `json:"volume"`
}
