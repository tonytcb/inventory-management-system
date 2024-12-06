CREATE TABLE IF NOT EXISTS transactions
(
    id           SERIAL PRIMARY KEY,
    reference_id INT                         NOT NULL,
    type         VARCHAR(50)                 NOT NULL,
    amount       DECIMAL(20, 8)              NOT NULL,
    fx_rate_id   INT                         NOT NULL,
    revenue      DECIMAL(20, 8)              NOT NULL,
    created_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITHOUT TIME ZONE,

    -- Foreign keys
    CONSTRAINT transactions_fx_rate_id_fk FOREIGN KEY (fx_rate_id) REFERENCES fx_rates (id),

    -- Unique constraint
    CONSTRAINT transactions_unique_reference_id_type UNIQUE (reference_id, type)
);

CREATE INDEX IF NOT EXISTS idx_transactions_reference_id ON transactions USING btree (reference_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions USING btree (type);
