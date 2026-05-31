package telemetry

import (
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseCollector(t *testing.T) {
	testDB := db.SetupTestDB(t)

	// Seed test users
	user1 := models.User{ID: 1, Username: "alice", CreatedAt: time.Now().AddDate(0, 0, -2)}
	user2 := models.User{ID: 2, Username: "bob", CreatedAt: time.Now().AddDate(0, 0, -10)}
	assert.NoError(t, testDB.Create(&user1).Error)
	assert.NoError(t, testDB.Create(&user2).Error)

	// Seed stats for leaderboard
	stats1 := models.UserStats{UserID: 1, Wins: 10, Losses: 2, Draws: 1, TotalGames: 13, AvgGameDurationSecs: 300}
	stats2 := models.UserStats{UserID: 2, Wins: 5, Losses: 7, Draws: 1, TotalGames: 13, AvgGameDurationSecs: 420}
	assert.NoError(t, testDB.Create(&stats1).Error)
	assert.NoError(t, testDB.Create(&stats2).Error)

	// Seed games & moves
	finishedTime := time.Now().Add(-5 * time.Minute)
	game1 := models.Game{
		ID:            "game-pvp-finished",
		Player1UserID: &user1.ID,
		Player2UserID: &user2.ID,
		GameType:      "pvp",
		InitialState:  "{}",
		CreatedAt:     time.Now().Add(-15 * time.Minute),
		FinishedAt:    &finishedTime,
	}
	game2 := models.Game{
		ID:            "game-ai-active",
		Player1UserID: &user1.ID,
		GameType:      "ai",
		InitialState:  "{}",
		CreatedAt:     time.Now().Add(-10 * time.Minute),
		FinishedAt:    nil,
	}
	assert.NoError(t, testDB.Create(&game1).Error)
	assert.NoError(t, testDB.Create(&game2).Error)

	move1 := models.GameMove{ID: 1, GameID: "game-pvp-finished", MoveIndex: 0, PlayerID: 1, FromX: 0, FromY: 0, ToX: 0, ToY: 1}
	move2 := models.GameMove{ID: 2, GameID: "game-pvp-finished", MoveIndex: 1, PlayerID: 2, FromX: 9, FromY: 9, ToX: 9, ToY: 8}
	assert.NoError(t, testDB.Create(&move1).Error)
	assert.NoError(t, testDB.Create(&move2).Error)

	collector := NewDatabaseCollector(0)

	t.Run("describe descriptors", func(t *testing.T) {
		descChan := make(chan *prometheus.Desc, 20)
		collector.Describe(descChan)
		close(descChan)

		var descriptors []*prometheus.Desc
		for desc := range descChan {
			descriptors = append(descriptors, desc)
		}
		assert.Len(t, descriptors, 13)
	})

	t.Run("collect metrics", func(t *testing.T) {
		reg := prometheus.NewRegistry()
		err := reg.Register(collector)
		assert.NoError(t, err)

		metricFamilies, err := reg.Gather()
		assert.NoError(t, err)

		findGauge := func(name string) float64 {
			for _, mf := range metricFamilies {
				if *mf.Name == name {
					return *mf.Metric[0].Gauge.Value
				}
			}
			return -1
		}

		assert.Equal(t, float64(2), findGauge("gostrategy_users_total"))
		assert.Equal(t, float64(1), findGauge("gostrategy_users_new_7d"))
		assert.Equal(t, float64(2), findGauge("gostrategy_games_total"))
		assert.Equal(t, float64(1), findGauge("gostrategy_games_active"))
		assert.Equal(t, float64(2), findGauge("gostrategy_moves_total"))

		// 10 minutes average duration in SQLite setup
		assert.InDelta(t, 10.0, findGauge("gostrategy_game_duration_avg_minutes"), 0.1)

		// Verify game type splits
		var gameTypesFound int
		for _, mf := range metricFamilies {
			if *mf.Name == "gostrategy_games_by_type" {
				for _, m := range mf.Metric {
					label := *m.Label[0].Value
					val := *m.Gauge.Value
					switch label {
					case "pvp":
						assert.Equal(t, float64(1), val)
						gameTypesFound++
					case "ai":
						assert.Equal(t, float64(1), val)
						gameTypesFound++
					}
				}
			}
		}
		assert.Equal(t, 2, gameTypesFound)

		// Verify leaderboard metrics
		var leaderboardAliceWins float64 = -1
		var leaderboardBobWins float64 = -1
		for _, mf := range metricFamilies {
			if *mf.Name == "gostrategy_user_wins" {
				for _, m := range mf.Metric {
					user := *m.Label[0].Value
					val := *m.Gauge.Value
					switch user {
					case "alice":
						leaderboardAliceWins = val
					case "bob":
						leaderboardBobWins = val
					}
				}
			}
		}
		assert.Equal(t, float64(10), leaderboardAliceWins)
		assert.Equal(t, float64(5), leaderboardBobWins)
	})
}
