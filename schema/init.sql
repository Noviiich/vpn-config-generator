CREATE TABLE IF NOT EXISTS users (
    telegram_id BIGINT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    private_key TEXT NOT NULL CHECK (length(private_key) > 10),
    public_key TEXT NOT NULL CHECK (length(private_key) > 10),
    ip_id INET UNIQUE NOT NULL REFERENCES ip_pool(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
);

CREATE TABLE IF NOT EXISTS ip_pool (
	id SERIAL PRIMARY KEY,
    ip inet UNIQUE NOT NULL,
    network CIDR NOT NULL UNIQUE,
    last_assigned_ip INET NOT NULL, 
);

CREATE TABLE IF NOT EXISTS subscription_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    duration INTERVAL NOT NULL,
    max_devices INT NOT NULL DEFAULT 3
);

CREATE TABLE IF NOT EXISTS subscriptions (
    user_id BIGINT PRIMARY KEY REFERENCES users(telegram_id) ON DELETE CASCADE,
    expiry_date TIMESTAMP DEFAULT NULL,
    is_active BOOLEAN DEFAULT FALSE
);