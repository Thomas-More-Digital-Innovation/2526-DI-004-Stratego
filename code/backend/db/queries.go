package db

import (
	"database/sql"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user with hashed password
func CreateUser(username, password, profilePicture string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var user models.User
	query := `
		INSERT INTO users (username, password_hash, profile_picture)
		VALUES ($1, $2, $3)
		RETURNING id, username, profile_picture, created_at, updated_at
	`
	err = DB.QueryRow(query, username, string(hashedPassword), profilePicture).Scan(
		&user.ID, &user.Username, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create initial stats for the user
	_, err = DB.Exec(`
		INSERT INTO user_stats (user_id)
		VALUES ($1)
	`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user stats: %w", err)
	}

	return &user, nil
}

// AuthenticateUser checks username and password, returns user if valid
func AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	var passwordHash string

	query := `
		SELECT id, username, password_hash, profile_picture, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	err := DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &passwordHash, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid username or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID int) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, username, profile_picture, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserStats retrieves stats for a user
func GetUserStats(userID int) (*models.UserStats, error) {
	var stats models.UserStats
	query := `
		SELECT id, user_id, total_games, wins, losses, draws, 
		       total_moves, avg_game_duration_seconds, created_at, updated_at
		FROM user_stats
		WHERE user_id = $1
	`
	err := DB.QueryRow(query, userID).Scan(
		&stats.ID, &stats.UserID, &stats.TotalGames, &stats.Wins, &stats.Losses,
		&stats.Draws, &stats.TotalMoves, &stats.AvgGameDurationSecs,
		&stats.CreatedAt, &stats.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	return &stats, nil
}

// UpdateUserStats updates game statistics for a user
func UpdateUserStats(userID int, won bool, moveCount int, durationSecs float64) error {
	query := `
		UPDATE user_stats
		SET total_games = total_games + 1,
		    wins = wins + $1,
		    losses = losses + $2,
		    total_moves = total_moves + $3,
		    avg_game_duration_seconds = (avg_game_duration_seconds * total_games + $4) / (total_games + 1)
		WHERE user_id = $5
	`
	winsInc := 0
	lossesInc := 0
	if won {
		winsInc = 1
	} else {
		lossesInc = 1
	}

	_, err := DB.Exec(query, winsInc, lossesInc, moveCount, durationSecs, userID)
	if err != nil {
		return fmt.Errorf("failed to update user stats: %w", err)
	}
	return nil
}

// CreateBoardSetup saves a new board setup
func CreateBoardSetup(userID int, name, description, setupData string, isDefault bool) (*models.BoardSetup, error) {
	var setup models.BoardSetup
	query := `
		INSERT INTO board_setups (user_id, name, description, setup_data, is_default)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, description, setup_data, is_default, created_at, updated_at
	`
	err := DB.QueryRow(query, userID, name, description, setupData, isDefault).Scan(
		&setup.ID, &setup.UserID, &setup.Name, &setup.Description,
		&setup.SetupData, &setup.IsDefault, &setup.CreatedAt, &setup.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create board setup: %w", err)
	}
	return &setup, nil
}

// GetBoardSetup retrieves a board setup by ID and verifying ownership
func GetBoardSetup(setupID, userID int) (*models.BoardSetup, error) {
	var setup models.BoardSetup
	query := `
		SELECT id, user_id, name, description, setup_data, is_default, created_at, updated_at
		FROM board_setups
		WHERE id = $1 AND user_id = $2
	`
	err := DB.QueryRow(query, setupID, userID).Scan(
		&setup.ID, &setup.UserID, &setup.Name, &setup.Description,
		&setup.SetupData, &setup.IsDefault, &setup.CreatedAt, &setup.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get board setup: %w", err)
	}
	return &setup, nil
}

// GetUserBoardSetups retrieves all board setups for a user
func GetUserBoardSetups(userID int) ([]models.BoardSetup, error) {
	query := `
		SELECT id, user_id, name, description, setup_data, is_default, created_at, updated_at
		FROM board_setups
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query board setups: %w", err)
	}
	defer rows.Close()

	var setups []models.BoardSetup
	for rows.Next() {
		var setup models.BoardSetup
		err := rows.Scan(
			&setup.ID, &setup.UserID, &setup.Name, &setup.Description,
			&setup.SetupData, &setup.IsDefault, &setup.CreatedAt, &setup.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan board setup: %w", err)
		}
		setups = append(setups, setup)
	}
	return setups, nil
}

// UpdateBoardSetup updates an existing board setup and verifying ownership
func UpdateBoardSetup(setupID, userID int, name, description, setupData string, isDefault bool) error {
	query := `
		UPDATE board_setups
		SET name = COALESCE(NULLIF($1, ''), name),
		    description = COALESCE(NULLIF($2, ''), description),
		    setup_data = COALESCE(NULLIF($3, ''), setup_data),
		    is_default = $4,
		    updated_at = $5
		WHERE id = $6 AND user_id = $7
	`
	result, err := DB.Exec(query, name, description, setupData, isDefault, time.Now(), setupID, userID)
	if err != nil {
		return fmt.Errorf("failed to update board setup: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("board setup not found or not owned by user")
	}

	return nil
}

// DeleteBoardSetup deletes a board setup
func DeleteBoardSetup(setupID, userID int) error {
	query := `DELETE FROM board_setups WHERE id = $1 AND user_id = $2`
	result, err := DB.Exec(query, setupID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete board setup: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("board setup not found or not owned by user")
	}

	return nil
}

// SaveGame persists the game metadata and initial state
func SaveGame(gameID string, p1ID, p2ID *int, gameType string, initialState interface{}, winnerID *int) error {
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return fmt.Errorf("failed to marshal initial state: %w", err)
	}

	query := `
		INSERT INTO games (id, player1_user_id, player2_user_id, winner_id, game_type, initial_state, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = DB.Exec(query, gameID, p1ID, p2ID, winnerID, gameType, stateJSON, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}
	return nil
}

// SaveMove persists a single move in a game's history
func SaveMove(gameID string, move models.HistoricalMove) error {
	attackerJSON, err := json.Marshal(move.Attacker)
	if err != nil {
		return fmt.Errorf("failed to marshal attacker data: %w", err)
	}
	defenderJSON, err := json.Marshal(move.Defender)
	if err != nil {
		return fmt.Errorf("failed to marshal defender data: %w", err)
	}

	query := `
		INSERT INTO game_moves (game_id, move_index, player_id, from_x, from_y, to_x, to_y, attacker_data, defender_data, result)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = DB.Exec(query, gameID, move.MoveIndex, move.PlayerID,
		move.FromX, move.FromY, move.ToX, move.ToY,
		attackerJSON, defenderJSON, move.Result)
	if err != nil {
		return fmt.Errorf("failed to save move: %w", err)
	}
	return nil
}
func GetGameHistory(gameID string) (*models.GameHistory, error) {
	var history models.GameHistory
	history.GameID = gameID

	var initialStateJSON []byte
	query := `
		SELECT initial_state, winner_id
		FROM games
		WHERE id = $1
	`
	err := DB.QueryRow(query, gameID).Scan(&initialStateJSON, &history.WinnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game history metadata: %w", err)
	}

	if err := json.Unmarshal(initialStateJSON, &history.InitialState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial state: %w", err)
	}

	query = `
		SELECT move_index, player_id, from_x, from_y, to_x, to_y, attacker_data, defender_data, result
		FROM game_moves
		WHERE game_id = $1
		ORDER BY move_index ASC
	`
	rows, err := DB.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to query game moves: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.HistoricalMove
		var attackerJSON, defenderJSON []byte
		err = rows.Scan(&m.MoveIndex, &m.PlayerID, &m.FromX, &m.FromY, &m.ToX, &m.ToY, &attackerJSON, &defenderJSON, &m.Result)
		if err != nil {
			return nil, fmt.Errorf("failed to scan historical move: %w", err)
		}

		if len(attackerJSON) > 0 {
			if err := json.Unmarshal(attackerJSON, &m.Attacker); err != nil {
				return nil, fmt.Errorf("failed to unmarshal attacker data: %w", err)
			}
		}
		if len(defenderJSON) > 0 {
			if err := json.Unmarshal(defenderJSON, &m.Defender); err != nil {
				return nil, fmt.Errorf("failed to unmarshal defender data: %w", err)
			}
		}

		history.Moves = append(history.Moves, m)
	}

	return &history, nil
}
