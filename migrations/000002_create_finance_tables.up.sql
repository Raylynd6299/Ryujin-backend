-- Migration: 000002_create_finance_tables
-- Creates all tables for the Finance bounded context.
-- Tables: categories, income_sources, expenses, debts, accounts

-- ============================================================
-- CATEGORIES
-- Stores both system-default categories (user_id IS NULL)
-- and user-created custom categories.
-- ============================================================
CREATE TABLE IF NOT EXISTS categories (
    id          UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID         NULL REFERENCES users(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(20)  NOT NULL CHECK (type IN ('income', 'expense', 'both')),
    icon        VARCHAR(50)  NOT NULL DEFAULT '',
    color       VARCHAR(20)  NOT NULL DEFAULT '',
    is_default  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_user_id ON categories (user_id);
CREATE INDEX idx_categories_type    ON categories (type);

-- ============================================================
-- INCOME SOURCES
-- Each row = one source of income (salary, freelance, etc.)
-- Amounts stored in integer cents.
-- ============================================================
CREATE TABLE IF NOT EXISTS income_sources (
    id           UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id  UUID         NULL REFERENCES categories(id) ON DELETE SET NULL,
    name         VARCHAR(150) NOT NULL,
    description  TEXT         NOT NULL DEFAULT '',
    amount_cents BIGINT       NOT NULL CHECK (amount_cents > 0),
    currency     VARCHAR(3)   NOT NULL DEFAULT 'USD',
    income_type  VARCHAR(30)  NOT NULL CHECK (income_type IN ('salary', 'freelance', 'rental', 'dividend', 'business', 'other')),
    recurrence   VARCHAR(20)  NOT NULL CHECK (recurrence IN ('none', 'daily', 'weekly', 'biweekly', 'monthly', 'quarterly', 'annually')),
    start_date   DATE         NOT NULL,
    end_date     DATE         NULL,
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_income_sources_user_id   ON income_sources (user_id);
CREATE INDEX idx_income_sources_is_active ON income_sources (user_id, is_active);

-- ============================================================
-- EXPENSES
-- Each row = one expense (rent, food, subscriptions, etc.)
-- Priority determines how necessary the expense is.
-- ============================================================
CREATE TABLE IF NOT EXISTS expenses (
    id           UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id  UUID         NULL REFERENCES categories(id) ON DELETE SET NULL,
    name         VARCHAR(150) NOT NULL,
    description  TEXT         NOT NULL DEFAULT '',
    amount_cents BIGINT       NOT NULL CHECK (amount_cents > 0),
    currency     VARCHAR(3)   NOT NULL DEFAULT 'USD',
    priority     VARCHAR(20)  NOT NULL CHECK (priority IN ('essential', 'important', 'optional', 'low')),
    recurrence   VARCHAR(20)  NOT NULL CHECK (recurrence IN ('none', 'daily', 'weekly', 'biweekly', 'monthly', 'quarterly', 'annually')),
    expense_date DATE         NOT NULL,
    end_date     DATE         NULL,
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expenses_user_id        ON expenses (user_id);
CREATE INDEX idx_expenses_is_active      ON expenses (user_id, is_active);
CREATE INDEX idx_expenses_expense_date   ON expenses (user_id, expense_date DESC);
CREATE INDEX idx_expenses_priority       ON expenses (user_id, priority);

-- ============================================================
-- DEBTS
-- Tracks financial liabilities (credit cards, loans, mortgages).
-- Amounts in integer cents. is_active = false when fully paid.
-- ============================================================
CREATE TABLE IF NOT EXISTS debts (
    id                     UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id                UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name                   VARCHAR(150) NOT NULL,
    description            TEXT         NOT NULL DEFAULT '',
    debt_type              VARCHAR(30)  NOT NULL CHECK (debt_type IN ('credit_card', 'personal_loan', 'mortgage', 'car_loan', 'student_loan', 'other')),
    total_amount_cents     BIGINT       NOT NULL CHECK (total_amount_cents > 0),
    remaining_amount_cents BIGINT       NOT NULL CHECK (remaining_amount_cents >= 0),
    monthly_payment_cents  BIGINT       NOT NULL CHECK (monthly_payment_cents > 0),
    currency               VARCHAR(3)   NOT NULL DEFAULT 'USD',
    interest_rate          NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (interest_rate >= 0),
    start_date             DATE         NULL,
    due_date               DATE         NULL,
    is_active              BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_remaining_lte_total CHECK (remaining_amount_cents <= total_amount_cents)
);

CREATE INDEX idx_debts_user_id   ON debts (user_id);
CREATE INDEX idx_debts_is_active ON debts (user_id, is_active);

-- ============================================================
-- ACCOUNTS
-- Tracks liquid financial accounts (bank accounts, cash, wallets).
-- Balance stored in integer cents.
-- ============================================================
CREATE TABLE IF NOT EXISTS accounts (
    id            UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id       UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          VARCHAR(150) NOT NULL,
    description   TEXT         NOT NULL DEFAULT '',
    account_type  VARCHAR(20)  NOT NULL CHECK (account_type IN ('checking', 'savings', 'cash', 'wallet')),
    balance_cents BIGINT       NOT NULL DEFAULT 0,
    currency      VARCHAR(3)   NOT NULL DEFAULT 'USD',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id   ON accounts (user_id);
CREATE INDEX idx_accounts_is_active ON accounts (user_id, is_active);

-- ============================================================
-- SEED: System default categories
-- These are global (user_id IS NULL) and cannot be deleted.
-- ============================================================
INSERT INTO categories (id, user_id, name, type, icon, color, is_default) VALUES
    -- Income categories
    (uuid_generate_v4(), NULL, 'Salary',        'income',  '💼', '#10B981', TRUE),
    (uuid_generate_v4(), NULL, 'Freelance',      'income',  '💻', '#6366F1', TRUE),
    (uuid_generate_v4(), NULL, 'Investments',    'income',  '📈', '#F59E0B', TRUE),
    (uuid_generate_v4(), NULL, 'Rental',         'income',  '🏠', '#8B5CF6', TRUE),
    (uuid_generate_v4(), NULL, 'Other Income',   'income',  '💰', '#64748B', TRUE),
    -- Expense categories
    (uuid_generate_v4(), NULL, 'Housing',        'expense', '🏠', '#EF4444', TRUE),
    (uuid_generate_v4(), NULL, 'Food',           'expense', '🍔', '#F97316', TRUE),
    (uuid_generate_v4(), NULL, 'Transport',      'expense', '🚗', '#3B82F6', TRUE),
    (uuid_generate_v4(), NULL, 'Health',         'expense', '🏥', '#EC4899', TRUE),
    (uuid_generate_v4(), NULL, 'Entertainment',  'expense', '🎬', '#A855F7', TRUE),
    (uuid_generate_v4(), NULL, 'Education',      'expense', '📚', '#14B8A6', TRUE),
    (uuid_generate_v4(), NULL, 'Subscriptions',  'expense', '📱', '#F59E0B', TRUE),
    (uuid_generate_v4(), NULL, 'Clothing',       'expense', '👕', '#6B7280', TRUE),
    (uuid_generate_v4(), NULL, 'Utilities',      'expense', '⚡', '#84CC16', TRUE),
    (uuid_generate_v4(), NULL, 'Other Expenses', 'expense', '📦', '#64748B', TRUE);
