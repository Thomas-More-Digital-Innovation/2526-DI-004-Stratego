package db

import (
	"database/sql"
	"digital-innovation/stratego/models"
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

// GetBoardSetup retrieves a board setup by ID
func GetBoardSetup(setupID int) (*models.BoardSetup, error) {
	var setup models.BoardSetup
	query := `
		SELECT id, user_id, name, description, setup_data, is_default, created_at, updated_at
		FROM board_setups
		WHERE id = $1
	`
	err := DB.QueryRow(query, setupID).Scan(
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

// UpdateBoardSetup updates an existing board setup
func UpdateBoardSetup(setupID int, name, description, setupData string, isDefault bool) error {
	query := `
		UPDATE board_setups
		SET name = COALESCE(NULLIF($1, ''), name),
		    description = COALESCE(NULLIF($2, ''), description),
		    setup_data = COALESCE(NULLIF($3, ''), setup_data),
		    is_default = $4,
		    updated_at = $5
		WHERE id = $6
	`
	_, err := DB.Exec(query, name, description, setupData, isDefault, time.Now(), setupID)
	if err != nil {
		return fmt.Errorf("failed to update board setup: %w", err)
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
