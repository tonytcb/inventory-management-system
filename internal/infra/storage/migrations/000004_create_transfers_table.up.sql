CREATE TABLE IF NOT EXISTS transfers
(
    id               SERIAL PRIMARY KEY,
    converted_amount DECIMAL(20, 8)              NOT NULL,
    final_amount     DECIMAL(20, 8)              NOT NULL,
    original_amount  DECIMAL(20, 8)              NOT NULL,
    description      TEXT,
    status           VARCHAR(50)                 NOT NULL,
    from_currency    VARCHAR(3)                  NOT NULL,
    to_currency      VARCHAR(3)                  NOT NULL,
    created_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITHOUT TIME ZONE
);
