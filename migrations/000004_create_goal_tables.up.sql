-- Migration: 000004_create_goal_tables
-- Creates all tables for the Goal bounded context.
-- Tables: purchase_goals, goal_contributions
-- ORDER MATTERS: purchase_goals must exist before goal_contributions (FK target).

-- ============================================================
-- PURCHASE GOALS
-- Each row = one savings goal (buy a car, laptop, vacation, etc.)
-- Target amount stored as BIGINT cents. Priority: low/medium/high.
-- is_completed flips to TRUE when total contributions >= target.
-- ============================================================
CREATE TABLE IF NOT EXISTS purchase_goals (
    id                   UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id              UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name                 VARCHAR(150) NOT NULL,
    description          TEXT         NOT NULL DEFAULT '',
    icon                 VARCHAR(10)  NOT NULL DEFAULT '',
    target_amount_cents  BIGINT       NOT NULL CHECK (target_amount_cents > 0),
    currency             VARCHAR(3)   NOT NULL DEFAULT 'USD',
    priority             VARCHAR(10)  NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    deadline             DATE         NULL,
    is_completed         BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_purchase_goals_user_id      ON purchase_goals (user_id);
CREATE INDEX IF NOT EXISTS idx_purchase_goals_is_completed ON purchase_goals (user_id, is_completed);
CREATE INDEX IF NOT EXISTS idx_purchase_goals_priority     ON purchase_goals (user_id, priority);
CREATE INDEX IF NOT EXISTS idx_purchase_goals_deadline     ON purchase_goals (user_id, deadline);

-- ============================================================
-- GOAL CONTRIBUTIONS
-- Records each deposit made toward a goal.
-- Amount stored as BIGINT cents. FK to purchase_goals with CASCADE.
-- user_id is denormalized here to make user-scoped queries fast.
-- ============================================================
CREATE TABLE IF NOT EXISTS goal_contributions (
    id           UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    goal_id      UUID         NOT NULL REFERENCES purchase_goals(id) ON DELETE CASCADE,
    user_id      UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount_cents BIGINT       NOT NULL CHECK (amount_cents > 0),
    currency     VARCHAR(3)   NOT NULL DEFAULT 'USD',
    date         DATE         NOT NULL,
    notes        VARCHAR(300) NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_goal_contributions_goal_id ON goal_contributions (goal_id);
CREATE INDEX IF NOT EXISTS idx_goal_contributions_user_id ON goal_contributions (user_id);
CREATE INDEX IF NOT EXISTS idx_goal_contributions_date    ON goal_contributions (goal_id, date DESC);
