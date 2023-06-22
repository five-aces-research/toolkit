-- name: ReadMinutes :many
SELECT starttime, open, high, close, low, volume
FROM klines
WHERE ticker_id = $1
  AND starttime > @st
  AND starttime < @et
ORDER BY starttime asc;

-- name: WriteMinutes :copyfrom
INSERT INTO klines (ticker_id, starttime,open,high,close,low,volume)
VALUES ($1,$2,$3,$4,$5,$6,$7);