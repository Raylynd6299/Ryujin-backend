-- Rollback: 000003_create_investment_tables
-- Drop in reverse dependency order: holdings first, then history, then quotes.

DROP TABLE IF EXISTS holdings;
DROP TABLE IF EXISTS stock_price_history;
DROP TABLE IF EXISTS stock_quotes;
