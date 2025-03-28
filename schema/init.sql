CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ip_pool (
	id SERIAL PRIMARY KEY,
    device_id INT REFERENCES devices(id) ON DELETE SET NULL,
    ip INET UNIQUE NOT NULL,
    is_available BOOLEAN GENERATED ALWAYS AS (device_id IS NULL) STORED
);

CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    private_key TEXT NOT NULL CHECK (length(private_key) > 10),
    public_key TEXT NOT NULL CHECK (length(public_key) > 10),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_active TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE TABLE IF NOT EXISTS subscription_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    duration INTERVAL NOT NULL,
    max_devices INT NOT NULL DEFAULT 3
);

CREATE TABLE IF NOT EXISTS subscriptions (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    type_id INT NOT NULL REFERENCES subscription_types(id),
    start_date TIMESTAMP NOT NULL DEFAULT NOW(),
    expiry_date TIMESTAMP DEFAULT NULL,
    is_active BOOLEAN DEFAULT FALSE
);