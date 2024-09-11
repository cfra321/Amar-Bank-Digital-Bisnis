-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    phonenumber VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(100),
    modified_at TIMESTAMP DEFAULT NOW(),
    modified_by VARCHAR(100),
    activated_by VARCHAR(100), -- Comma added here
    activated_at TIMESTAMP   DEFAULT NOW()    -- Added type for activated_at
)

-- +migrate StatementEnd
