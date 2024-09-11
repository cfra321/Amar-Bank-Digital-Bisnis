-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE transaction_logs (
    id SERIAL PRIMARY KEY,
    transfer_id INT NOT NULL REFERENCES transfers(id),
    log_message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +migrate StatementEnd
