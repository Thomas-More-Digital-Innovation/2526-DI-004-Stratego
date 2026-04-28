-- Dev only
-- DROP TABLE IF EXISTS board_setups CASCADE;
-- DROP TABLE IF EXISTS user_stats CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  profile_picture VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_stats (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) UNIQUE NOT NULL,
  total_games INTEGER NOT NULL DEFAULT 0,
  wins INTEGER NOT NULL DEFAULT 0,
  losses INTEGER NOT NULL DEFAULT 0,
  draws INTEGER NOT NULL DEFAULT 0,
  total_moves INTEGER NOT NULL DEFAULT 0,
  avg_game_duration_seconds REAL NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS board_setups (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  setup_data VARCHAR(40) NOT NULL,
  is_default BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_board_setups_user_id ON board_setups(user_id);
CREATE INDEX idx_user_stats_user_id ON user_stats(user_id);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER user_stats_updated_at
BEFORE UPDATE ON user_stats
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER board_setups_updated_at
BEFORE UPDATE ON board_setups
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TABLE IF NOT EXISTS games (
  id VARCHAR(100) PRIMARY KEY,
  player1_user_id INTEGER REFERENCES users(id),
  player2_user_id INTEGER REFERENCES users(id),
  winner_id INTEGER, -- 0 for Player 1, 1 for Player 2, NULL for draw
  game_type VARCHAR(50) NOT NULL,
  initial_state JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  finished_at TIMESTAMPTZ -- NULL if abandoned/in progress
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
  attacker_data JSONB, -- Optional combat data (rank, type)
  defender_data JSONB, -- Optional combat data (rank, type)
  result VARCHAR(20) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_game_moves_game_id ON game_moves(game_id);
CREATE INDEX idx_games_player1_id ON games(player1_user_id);
CREATE INDEX idx_games_player2_id ON games(player2_user_id);
