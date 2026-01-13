CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    nickname TEXT,
    password TEXT NOT NULL,
    photo TEXT,
    city TEXT,
    status TEXT,
    agreement_pd BOOlEAN,
    agreement_ea BOOlEAN
);