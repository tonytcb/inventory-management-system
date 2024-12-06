CREATE TABLE IF NOT EXISTS currency_pools_ledger
(
    id            SERIAL PRIMARY KEY,
    currency_code VARCHAR(3)                                            NOT NULL,
    balance       DECIMAL(20, 8)                                        NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_currency_pools_ledger_currency_code ON currency_pools_ledger USING btree (currency_code);

-- Populate the table with initial data
INSERT INTO currency_pools_ledger (currency_code, balance)
VALUES ('USD', '1000000'),
       ('EUR', '921658'),
       ('JPY', '109890110'),
       ('GBP', '750000'),
       ('AUD', '1349528');
