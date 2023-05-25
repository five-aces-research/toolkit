-- name: CreateTicker :one
INSERT INTO tickers(exchange_id, name) VALUES ($1,$2) RETURNING ticker_id;

-- name: GetTickerId :one
SELECT ticker_id FROM tickers WHERE exchange_id = $1 and name = $2;


-- name: ReadOHCL :many
SELECT starttime, open, high, close, low, volume
FROM klines
WHERE ticker_id = $1
  AND resolution = $2
  AND starttime > @st
  AND starttime < @et
ORDER BY starttime asc;


-- name: WriteOHCLV :copyfrom
INSERT INTO klines (ticker_id, resolution, starttime,open,high,close,low,volume)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8);


-- name: GetTickers :many
SELECT name FROM tickers;

-- name: MinAndMax :one
SELECT MAX(starttime)::timestamp as Max, Min(starttime)::timestamp as Min FROM klines
WHERE ticker_id = $1 and resolution = $2;