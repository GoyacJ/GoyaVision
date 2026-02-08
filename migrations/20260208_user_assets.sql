-- 用户资产管理相关表迁移脚本

-- 1. 用户余额与积分表
CREATE TABLE IF NOT EXISTS user_balances (
    user_id UUID PRIMARY KEY,
    balance DECIMAL(16,2) DEFAULT 0,
    points BIGINT DEFAULT 0,
    level VARCHAR(32) DEFAULT 'Free',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_user_balances_deleted_at ON user_balances(deleted_at);

-- 2. 用户订阅表
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    plan_name VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_deleted_at ON user_subscriptions(deleted_at);

-- 3. 交易记录表
CREATE TABLE IF NOT EXISTS transaction_records (
    id VARCHAR(64) PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(32) NOT NULL,
    method VARCHAR(32) NOT NULL,
    amount DECIMAL(16,2) NOT NULL,
    status VARCHAR(32) NOT NULL,
    remark TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_transaction_records_user_id ON transaction_records(user_id);
CREATE INDEX IF NOT EXISTS idx_transaction_records_deleted_at ON transaction_records(deleted_at);

-- 4. 积分记录表
CREATE TABLE IF NOT EXISTS point_records (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(32) NOT NULL,
    change BIGINT NOT NULL,
    balance BIGINT NOT NULL,
    remark TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_point_records_user_id ON point_records(user_id);

-- 5. 使用统计表
CREATE TABLE IF NOT EXISTS usage_stats (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    operator_calls BIGINT DEFAULT 0,
    ai_model_calls BIGINT DEFAULT 0,
    token_usage BIGINT DEFAULT 0,
    date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_usage_stats_user_id ON usage_stats(user_id);
CREATE INDEX IF NOT EXISTS idx_usage_stats_date ON usage_stats(date);
CREATE INDEX IF NOT EXISTS idx_usage_stats_deleted_at ON usage_stats(deleted_at);
