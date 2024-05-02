CREATE DATABASE userAuth;
\c userAuth;

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    verification_code INT NOT NULL CHECK (verification_code >= 100000 AND verification_code <= 999999),
    sent_at TIMESTAMPTZ NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE
);