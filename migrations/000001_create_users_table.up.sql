-- Migration: 000001_create_users_table
-- Creates the users table for the User bounded context.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    -- Identity
    id                           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Authentication
    email                        VARCHAR(255) NOT NULL UNIQUE,
    hashed_password              TEXT         NOT NULL,

    -- Profile
    first_name                   VARCHAR(100) NOT NULL,
    last_name                    VARCHAR(100) NOT NULL,

    -- Preferences
    default_savings_currency     VARCHAR(3)   NOT NULL DEFAULT 'USD',
    default_investment_currency  VARCHAR(3)   NOT NULL DEFAULT 'USD',
    locale                       VARCHAR(10)  NOT NULL DEFAULT 'en',

    -- Timestamps
    created_at                   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at                   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at                   TIMESTAMPTZ  NULL
);

-- Index for email lookups (login)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email) WHERE deleted_at IS NULL;

-- Index for soft-delete queries
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
