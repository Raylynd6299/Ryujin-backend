-- Migration: 000002_create_finance_tables (rollback)
-- Drop all finance tables in reverse order (respects FK dependencies)

DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS debts;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS income_sources;
DROP TABLE IF EXISTS categories;
