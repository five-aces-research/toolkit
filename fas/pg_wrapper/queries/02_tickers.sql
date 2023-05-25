-- name: GetAvaibleTickers :many
SELECT exchanges.name as exchange, tickers.name as ticker, tickers.ticker_id FROM tickers
JOIN exchanges on exchanges.exchange_id = tickers.exchange_id;


