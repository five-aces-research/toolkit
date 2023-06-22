CREATE TABLE IF NOT EXISTS exchanges(
    exchange_id SERIAL PRIMARY KEY,
    name text UNIQUE NOT NULL
);

CREATE TABLE if not exists tickers (
    ticker_id SERIAL PRIMARY KEY,
    exchange_id INT references exchanges(exchange_id) NOT NULL ,
    name text NOT NULL ,
    unique(exchange_id,name)
);

CREATE TABLE IF NOT EXISTS klines(
    ticker_id int references tickers(ticker_id) NOT NULL,
    resolution int       NOT NULL,
    starttime  BIGINT NOT NULL,
    open       float8    not null,
    high       float8    not null,
    close      float8    not null,
    low        float8    not null,
    volume     float8    not null,
    unique (ticker_id, resolution, starttime)
);


CREATE TABLE IF NOT EXISTS minutes(
    ticker_id int references tickers(ticker_id) NOT NULL,
    starttime  BIGINT NOT NULL,
    open       float8    not null,
    high       float8    not null,
    close      float8    not null,
    low        float8    not null,
    volume     float8    not null,
    unique (ticker_id, starttime)
);




INSERT INTO exchanges (name) VALUES('BYBIT');
INSERT INTO exchanges (name) VALUES('DERIBIT');

---- create above / drop below ----

DROP TABLE IF EXISTS minutes;
DROP TABLE IF EXISTS klines;
DROP TABLE IF EXISTS tickers;
DROP TABLE IF EXISTS exchanges;
