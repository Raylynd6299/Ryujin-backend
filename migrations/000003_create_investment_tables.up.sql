-- Migration: 000003_create_investment_tables
-- Creates all tables for the Investment bounded context.
-- Tables: stock_quotes, stock_price_history, holdings
-- ORDER MATTERS: stock_quotes must exist before holdings (FK target).

-- ============================================================
-- STOCK QUOTES
-- Global (non-user-scoped) cache of market data per ticker.
-- symbol is the natural PK (uppercase ticker, max 10 chars).
-- All monetary fields stored as BIGINT cents.
-- ============================================================
CREATE TABLE IF NOT EXISTS stock_quotes (
    symbol               VARCHAR(10)  PRIMARY KEY,
    name                 VARCHAR(150) NOT NULL,
    currency             VARCHAR(3)   NOT NULL,

    -- Price fields — all in cents (smallest currency unit)
    price_cents          BIGINT       NOT NULL DEFAULT 0,
    previous_close_cents BIGINT       NOT NULL DEFAULT 0,
    open_cents           BIGINT       NOT NULL DEFAULT 0,
    day_high_cents       BIGINT       NOT NULL DEFAULT 0,
    day_low_cents        BIGINT       NOT NULL DEFAULT 0,
    volume               BIGINT       NOT NULL DEFAULT 0,
    market_cap_cents     BIGINT       NOT NULL DEFAULT 0,
    week52_high_cents    BIGINT       NOT NULL DEFAULT 0,
    week52_low_cents     BIGINT       NOT NULL DEFAULT 0,

    -- Dimensionless ratios — float is acceptable here
    trailing_pe          NUMERIC(10,4) NOT NULL DEFAULT 0,
    forward_pe           NUMERIC(10,4) NOT NULL DEFAULT 0,
    dividend_yield       NUMERIC(10,4) NOT NULL DEFAULT 0,
    eps                  NUMERIC(10,4) NOT NULL DEFAULT 0,

    fetched_at           TIMESTAMPTZ  NOT NULL,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================================
-- STOCK PRICE HISTORY
-- Append-only snapshots of prices per symbol over time.
-- No updates or deletes — purely additive.
-- ============================================================
CREATE TABLE IF NOT EXISTS stock_price_history (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol      VARCHAR(10) NOT NULL REFERENCES stock_quotes(symbol) ON DELETE CASCADE,
    price_cents BIGINT      NOT NULL,
    currency    VARCHAR(3)  NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stock_price_history_symbol      ON stock_price_history (symbol);
CREATE INDEX IF NOT EXISTS idx_stock_price_history_recorded_at ON stock_price_history (symbol, recorded_at DESC);

-- ============================================================
-- HOLDINGS
-- User investment positions. symbol is a FK to stock_quotes
-- so every holding is backed by a real market data entry.
-- ============================================================
CREATE TABLE IF NOT EXISTS holdings (
    id                  UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- FK to stock_quotes — must exist before inserting a holding
    symbol              VARCHAR(10)  NOT NULL REFERENCES stock_quotes(symbol) ON DELETE RESTRICT,

    name                VARCHAR(150) NOT NULL,
    asset_type          VARCHAR(20)  NOT NULL CHECK (asset_type IN ('stock', 'etf', 'fixed_income', 'crypto', 'other')),

    -- Quantity stored as micro-units (1 share = 1_000_000) to avoid decimals
    quantity_micro      BIGINT       NOT NULL CHECK (quantity_micro > 0),

    -- Buy price per unit in cents
    buy_price_cents     BIGINT       NOT NULL CHECK (buy_price_cents > 0),
    buy_currency        VARCHAR(3)   NOT NULL DEFAULT 'USD',

    -- Current market price — NULL until first price refresh
    current_price_cents BIGINT       NULL,
    priced_at           TIMESTAMPTZ  NULL,

    notes               TEXT         NOT NULL DEFAULT '',

    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_holdings_user_id     ON holdings (user_id);
CREATE INDEX IF NOT EXISTS idx_holdings_symbol       ON holdings (symbol);
CREATE INDEX IF NOT EXISTS idx_holdings_user_symbol  ON holdings (user_id, symbol);
