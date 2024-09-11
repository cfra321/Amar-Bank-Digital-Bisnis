-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    account_number VARCHAR(20) NOT NULL UNIQUE,
    balance DECIMAL(18, 2) DEFAULT 0,
    account_type VARCHAR(10) NOT NULL,  -- 'internal' or 'external'
    created_at TIMESTAMP DEFAULT NOW()
);

-- +migrate StatementEnd
