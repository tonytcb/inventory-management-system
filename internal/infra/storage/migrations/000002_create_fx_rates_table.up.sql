CREATE TABLE IF NOT EXISTS fx_rates
(
    id            SERIAL PRIMARY KEY,
    from_currency VARCHAR(3)                  NOT NULL,
    to_currency   VARCHAR(3)                  NOT NULL,
    rate          DECIMAL(20, 8)              NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fx_rates_from_to_currency ON fx_rates (from_currency, to_currency);
