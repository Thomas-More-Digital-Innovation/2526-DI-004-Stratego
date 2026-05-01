-- Add games and game_moves tables for history tracking
CREATE TABLE IF NOT EXISTS games (
  id VARCHAR(100) PRIMARY KEY,
  player1_user_id INTEGER REFERENCES users(id),
  player2_user_id INTEGER REFERENCES users(id),
  winner_id INTEGER,
  game_type VARCHAR(50) NOT NULL,
  initial_state JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  finished_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS game_moves (
  id SERIAL PRIMARY KEY,
  game_id VARCHAR(100) REFERENCES games(id) ON DELETE CASCADE,
  move_index INTEGER NOT NULL,
  player_id INTEGER NOT NULL,
  from_x INTEGER NOT NULL,
  from_y INTEGER NOT NULL,
  to_x INTEGER NOT NULL,
  to_y INTEGER NOT NULL,
  attacker_data JSONB,
  defender_data JSONB,
  result VARCHAR(20) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_game_moves_game_id ON game_moves(game_id);
CREATE INDEX IF NOT EXISTS idx_games_player1_id ON games(player1_user_id);
CREATE INDEX IF NOT EXISTS idx_games_player2_id ON games(player2_user_id);
