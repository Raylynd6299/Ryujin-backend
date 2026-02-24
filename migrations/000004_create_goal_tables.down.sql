-- Rollback: 000004_create_goal_tables
-- Drops all tables created by the up migration in reverse dependency order.

DROP TABLE IF EXISTS goal_contributions;
DROP TABLE IF EXISTS purchase_goals;
