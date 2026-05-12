-- 005_games_rls.sql

-- Enable Row Level Security (RLS) on games and game_moves
ALTER TABLE games ENABLE ROW LEVEL SECURITY;
ALTER TABLE games FORCE ROW LEVEL SECURITY;

ALTER TABLE game_moves ENABLE ROW LEVEL SECURITY;
ALTER TABLE game_moves FORCE ROW LEVEL SECURITY;

-- Games: Anyone can view game metadata/status (Public Spectating)
DROP POLICY IF EXISTS games_isolation_policy ON games;
CREATE POLICY games_public_select_policy ON games
    FOR SELECT
    USING (true);

-- Game Moves: Anyone can view move history for any game
DROP POLICY IF EXISTS game_moves_isolation_policy ON game_moves;
CREATE POLICY game_moves_public_select_policy ON game_moves
    FOR SELECT
    USING (true);

-- Allow inserting games (Server creates games)
DROP POLICY IF EXISTS games_insert_policy ON games;
CREATE POLICY games_insert_policy ON games
    FOR INSERT
    WITH CHECK (true); 

-- Allow inserting moves (Server creates moves)
DROP POLICY IF EXISTS game_moves_insert_policy ON game_moves;
CREATE POLICY game_moves_insert_policy ON game_moves
    FOR INSERT
    WITH CHECK (true);

-- Note: No UPDATE or DELETE policies are created for users, 
-- ensuring spectators can view but not alter data.
