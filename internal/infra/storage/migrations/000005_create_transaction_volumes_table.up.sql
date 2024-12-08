CREATE TABLE IF NOT EXISTS transaction_volumes
(
    from_currency VARCHAR(3)                  NOT NULL,
    to_currency   VARCHAR(3)                  NOT NULL,
    volume        DECIMAL(20, 8)              NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL,

    PRIMARY KEY (from_currency, to_currency)
);
