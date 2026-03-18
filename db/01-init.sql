CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    nickname TEXT NOT NULL,
    password TEXT NOT NULL,
    user_role TEXT NOT NULL CHECK (user_role IN ('user', 'moderator', 'admin')),
    photo TEXT,
    city TEXT,
    status TEXT,
    points INTEGER,
    agreement_pd BOOlEAN,
    agreement_ea BOOlEAN,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS email_verifications (
    id BIGSERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    code TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NOW() + interval '10 minutes'
);