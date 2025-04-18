CREATE TABLE IF NOT EXISTS users (
    telegram_id BIGINT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    subscription_active BOOLEAN DEFAULT FALSE,
    subscription_expiry TIMESTAMP
);

CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(telegram_id) ON DELETE CASCADE,
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL,
    ip TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS ip_pool (
	id SERIAL PRIMARY KEY,
	last_ip TEXT NOT NULL
);