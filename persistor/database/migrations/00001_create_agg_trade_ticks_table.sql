-- +goose Up
-- +goose StatementBegin
CREATE TABLE agg_trade_ticks (
    symbol text not null,
    timestamp timestamp with time zone not null,
    open double precision not null,
    high double precision not null,
    low double precision not null,
    close double precision not null,
    volume double precision not null,
    PRIMARY KEY (symbol, timestamp)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE agg_trade_ticks;
-- +goose StatementEnd
