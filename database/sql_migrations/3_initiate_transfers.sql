-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE transfers (
    id SERIAL PRIMARY KEY,
    sender_account_id INT NOT NULL REFERENCES accounts(id),
    receiver_account_id INT NOT NULL REFERENCES accounts(id),
    amount DECIMAL(18, 2) NOT NULL,
    transfer_type VARCHAR(10) NOT NULL,  -- 'overbook' or 'bifast'
    fee DECIMAL(18, 2) DEFAULT 0,
    status VARCHAR(20) NOT NULL,  -- 'pending', 'completed', or 'failed'
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- +migrate StatementEnd
