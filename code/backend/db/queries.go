package db

import (
	"context"
	"digital-innovation/stratego/models"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser creates a new user with hashed password
func CreateUser(ctx context.Context, username, password, profilePicture string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Username:       username,
		PasswordHash:   string(hashedPassword),
		ProfilePicture: profilePicture,
	}

	err = DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if tx.Dialector.Name() == "postgres" {
			if err := tx.Exec(fmt.Sprintf("SET LOCAL app.current_user_id = '%d'", user.ID)).Error; err != nil {
				return err
			}
		}
		stats := models.UserStats{UserID: user.ID}
		return tx.Create(&stats).Error
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// UpdateUserPassword updates a user's password with a new hash
func UpdateUserPassword(ctx context.Context, userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password_hash", string(hashedPassword)).Error
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// AuthenticateUser checks username and password, returns user if valid
func AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User

	err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid username or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := DB.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserStats retrieves stats for a user
func GetUserStats(ctx context.Context, userID int) (*models.UserStats, error) {
	var stats models.UserStats
	err := DB.WithContext(ctx).Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	return &stats, nil
}

// UpdateUserStats updates game statistics for a user
func UpdateUserStats(ctx context.Context, userID int, won bool, moveCount int, durationSecs float64) error {
	winsIncrement, lossesIncrement := 0, 0
	if won {
		winsIncrement = 1
	} else {
		lossesIncrement = 1
	}

	err := DB.WithContext(ctx).Model(&models.UserStats{}).Where("user_id = ?", userID).
		Updates(map[string]any{
			"total_games":               gorm.Expr("total_games + 1"),
			"wins":                      gorm.Expr("wins + ?", winsIncrement),
			"losses":                    gorm.Expr("losses + ?", lossesIncrement),
			"total_moves":               gorm.Expr("total_moves + ?", moveCount),
			"avg_game_duration_seconds": gorm.Expr("(avg_game_duration_seconds * total_games + ?) / (total_games + 1)", durationSecs),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to update user stats: %w", err)
	}
	return nil
}

// CreateBoardSetup saves a new board setup
func CreateBoardSetup(ctx context.Context, userID int, name, description, setupData string, isDefault bool) (*models.BoardSetup, error) {
	setup := models.BoardSetup{
		UserID:      userID,
		Name:        name,
		Description: description,
		SetupData:   setupData,
		IsDefault:   isDefault,
	}
	err := WithRLS(ctx, func(tx *gorm.DB) error {
		return tx.Create(&setup).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create board setup: %w", err)
	}
	return &setup, nil
}

// GetBoardSetup retrieves a board setup by ID and verifying ownership
func GetBoardSetup(ctx context.Context, setupID, userID int) (*models.BoardSetup, error) {
	var setup models.BoardSetup
	err := DB.WithContext(ctx).Where("id = ? AND user_id = ?", setupID, userID).First(&setup).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get board setup: %w", err)
	}
	return &setup, nil
}

// GetUserBoardSetups retrieves all board setups for a user
func GetUserBoardSetups(ctx context.Context, userID int) ([]models.BoardSetup, error) {
	var setups []models.BoardSetup
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Order("is_default desc, created_at desc").Find(&setups).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query board setups: %w", err)
	}
	return setups, nil
}

// UpdateBoardSetup updates an existing board setup and verifying ownership
func UpdateBoardSetup(ctx context.Context, setupID, userID int, name, description, setupData string, isDefault bool) error {
	updates := map[string]any{
		"is_default": isDefault,
	}
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if setupData != "" {
		updates["setup_data"] = setupData
	}

	result := DB.WithContext(ctx).Model(&models.BoardSetup{}).Where("id = ? AND user_id = ?", setupID, userID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update board setup: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("board setup not found or not owned by user")
	}
	return nil
}

// DeleteBoardSetup deletes a board setup
func DeleteBoardSetup(ctx context.Context, setupID, userID int) error {
	result := DB.WithContext(ctx).Where("id = ? AND user_id = ?", setupID, userID).Delete(&models.BoardSetup{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete board setup: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("board setup not found or not owned by user")
	}
	return nil
}

// SaveGame persists the game metadata and initial state
func SaveGame(ctx context.Context, gameID string, p1ID, p2ID *int, gameType string, initialState interface{}, winnerID *int) error {
	stateJSON, err := json.Marshal(initialState)
	if err != nil {
		return fmt.Errorf("failed to marshal initial state: %w", err)
	}

	game := models.Game{
		ID:            gameID,
		Player1UserID: p1ID,
		Player2UserID: p2ID,
		WinnerID:      winnerID,
		GameType:      gameType,
		InitialState:  string(stateJSON),
	}

	err = DB.WithContext(ctx).Create(&game).Error
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}
	return nil
}

// SaveMove persists a single move in a game's history
func SaveMove(ctx context.Context, gameID string, move models.HistoricalMove) error {
	var attackerData, defenderData string

	if move.Attacker != nil {
		b, _ := json.Marshal(move.Attacker)
		attackerData = string(b)
	}
	if move.Defender != nil {
		b, _ := json.Marshal(move.Defender)
		defenderData = string(b)
	}

	gameMove := models.GameMove{
		GameID:       gameID,
		MoveIndex:    move.MoveIndex,
		PlayerID:     move.PlayerID,
		FromX:        move.FromX,
		FromY:        move.FromY,
		ToX:          move.ToX,
		ToY:          move.ToY,
		AttackerData: attackerData,
		DefenderData: defenderData,
		Result:       move.Result,
	}

	err := DB.WithContext(ctx).Create(&gameMove).Error
	if err != nil {
		return fmt.Errorf("failed to save move: %w", err)
	}
	return nil
}

// SaveMoves persists multiple moves in a single transaction
func SaveMoves(ctx context.Context, gameID string, moves []models.HistoricalMove) error {
	if len(moves) == 0 {
		return nil
	}

	var gameMoves []models.GameMove
	for _, move := range moves {
		var attackerData, defenderData string

		if move.Attacker != nil {
			b, _ := json.Marshal(move.Attacker)
			attackerData = string(b)
		}
		if move.Defender != nil {
			b, _ := json.Marshal(move.Defender)
			defenderData = string(b)
		}

		gameMoves = append(gameMoves, models.GameMove{
			GameID:       gameID,
			MoveIndex:    move.MoveIndex,
			PlayerID:     move.PlayerID,
			FromX:        move.FromX,
			FromY:        move.FromY,
			ToX:          move.ToX,
			ToY:          move.ToY,
			AttackerData: attackerData,
			DefenderData: defenderData,
			Result:       move.Result,
		})
	}

	err := DB.WithContext(ctx).Create(&gameMoves).Error
	if err != nil {
		return fmt.Errorf("failed to save moves: %w", err)
	}
	return nil
}

func GetGameHistory(ctx context.Context, gameID string) (*models.GameHistory, error) {
	var game models.Game
	err := DB.WithContext(ctx).Preload("Moves", func(db *gorm.DB) *gorm.DB {
		return db.Order("move_index ASC")
	}).Where("id = ?", gameID).First(&game).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get game history: %w", err)
	}

	var history models.GameHistory
	history.GameID = game.ID
	history.WinnerID = game.WinnerID

	if err := json.Unmarshal([]byte(game.InitialState), &history.InitialState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial state: %w", err)
	}

	for _, gm := range game.Moves {
		m := models.HistoricalMove{
			MoveIndex: gm.MoveIndex,
			PlayerID:  gm.PlayerID,
			FromX:     gm.FromX,
			FromY:     gm.FromY,
			ToX:       gm.ToX,
			ToY:       gm.ToY,
			Result:    gm.Result,
		}

		if gm.AttackerData != "" {
			var attacker models.PieceData
			if err := json.Unmarshal([]byte(gm.AttackerData), &attacker); err == nil {
				m.Attacker = &attacker
			}
		}
		if gm.DefenderData != "" {
			var defender models.PieceData
			if err := json.Unmarshal([]byte(gm.DefenderData), &defender); err == nil {
				m.Defender = &defender
			}
		}

		history.Moves = append(history.Moves, m)
	}

	return &history, nil
}

// SaveRefreshToken stores a new refresh token for a user
func SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	rt := models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	err := DB.WithContext(ctx).Create(&rt).Error
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

// GetUserIDByRefreshToken validates a refresh token and returns the owner's ID
func GetUserIDByRefreshToken(ctx context.Context, token string) (int, error) {
	var rt models.RefreshToken
	err := DB.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&rt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("invalid or expired refresh token")
		}
		return 0, fmt.Errorf("database error: %w", err)
	}
	return rt.UserID, nil
}

// DeleteRefreshToken removes a refresh token (e.g., on logout)
func DeleteRefreshToken(ctx context.Context, token string) error {
	return DB.WithContext(ctx).Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

// DeleteAllUserRefreshTokens revokes all sessions for a user
func DeleteAllUserRefreshTokens(ctx context.Context, userID int) error {
	return DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
