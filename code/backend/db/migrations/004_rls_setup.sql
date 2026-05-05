-- 004_rls_setup.sql

-- Enable Row Level Security (RLS) on relevant tables
ALTER TABLE board_setups ENABLE ROW LEVEL SECURITY;
ALTER TABLE refresh_tokens ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_stats ENABLE ROW LEVEL SECURITY;

-- Board Setups: Users can only see and modify their own setups
-- Using 'app.current_user_id' custom setting to enforce this
CREATE POLICY board_setups_isolation_policy ON board_setups
    USING (user_id = current_setting('app.current_user_id', true)::integer)
    WITH CHECK (user_id = current_setting('app.current_user_id', true)::integer);

-- Refresh Tokens: Users can only see and modify their own refresh tokens
CREATE POLICY refresh_tokens_isolation_policy ON refresh_tokens
    USING (user_id = current_setting('app.current_user_id', true)::integer)
    WITH CHECK (user_id = current_setting('app.current_user_id', true)::integer);

-- User Stats: Publicly viewable, but only owner can update
CREATE POLICY user_stats_select_policy ON user_stats
    FOR SELECT USING (true);

CREATE POLICY user_stats_modify_policy ON user_stats
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true)::integer)
    WITH CHECK (user_id = current_setting('app.current_user_id', true)::integer);
