CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    nickname TEXT NOT NULL,
    password TEXT NOT NULL,
    user_role TEXT NOT NULL CHECK (user_role IN ('user', 'moderator', 'admin')),
    photo TEXT,
    city TEXT,
    status TEXT,
    agreement_pd BOOlEAN,
    agreement_ea BOOlEAN
);