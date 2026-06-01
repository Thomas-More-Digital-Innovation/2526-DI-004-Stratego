// Package telemetry provides metrics collection and telemetry filtering capabilities.
package telemetry

import (
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/models"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelType   = "game_type"
	labelStatus = "status"
	labelUser   = "username"
)

type dbMetrics struct {
	usersTotal          float64
	usersNew7d          float64
	gamesTotal          float64
	gamesActive         float64
	gameDurationAvgMins float64
	gamesByType         map[string]float64
	gamesByStatus       map[string]float64
	movesTotal          float64
	leaderboard         []leaderboardEntry
}

type leaderboardEntry struct {
	username           string
	wins               float64
	losses             float64
	draws              float64
	games              float64
	avgDurationMinutes float64
}

// DatabaseCollector implements prometheus.Collector with an in-memory cache.
type DatabaseCollector struct {
	mu          sync.RWMutex
	cache       dbMetrics
	lastUpdated time.Time
	cacheTTL    time.Duration

	// Descriptions
	usersTotalDesc          *prometheus.Desc
	usersNew7dDesc          *prometheus.Desc
	gamesTotalDesc          *prometheus.Desc
	gamesActiveDesc         *prometheus.Desc
	gameDurationAvgMinsDesc *prometheus.Desc
	gamesByTypeDesc         *prometheus.Desc
	gamesByStatusDesc       *prometheus.Desc
	movesTotalDesc          *prometheus.Desc
	userWinsDesc            *prometheus.Desc
	userLossesDesc          *prometheus.Desc
	userDrawsDesc           *prometheus.Desc
	userGamesDesc           *prometheus.Desc
	userAvgDurationDesc     *prometheus.Desc
}

// NewDatabaseCollector creates a new DatabaseCollector with specified cache TTL.
func NewDatabaseCollector(ttl time.Duration) *DatabaseCollector {
	return &DatabaseCollector{
		cacheTTL: ttl,
		usersTotalDesc: prometheus.NewDesc(
			"gostrategy_users_total",
			"Total registered users",
			nil, nil,
		),
		usersNew7dDesc: prometheus.NewDesc(
			"gostrategy_users_new_7d",
			"Users registered in the last 7 days",
			nil, nil,
		),
		gamesTotalDesc: prometheus.NewDesc(
			"gostrategy_games_total",
			"Total games played",
			nil, nil,
		),
		gamesActiveDesc: prometheus.NewDesc(
			"gostrategy_games_active",
			"Active ongoing games",
			nil, nil,
		),
		gameDurationAvgMinsDesc: prometheus.NewDesc(
			"gostrategy_game_duration_avg_minutes",
			"Average game duration in minutes",
			nil, nil,
		),
		gamesByTypeDesc: prometheus.NewDesc(
			"gostrategy_games_by_type",
			"Total games by game type",
			[]string{labelType}, nil,
		),
		gamesByStatusDesc: prometheus.NewDesc(
			"gostrategy_games_by_status",
			"Games finished vs aborted/active",
			[]string{labelStatus}, nil,
		),
		movesTotalDesc: prometheus.NewDesc(
			"gostrategy_moves_total",
			"Total game moves played",
			nil, nil,
		),
		userWinsDesc: prometheus.NewDesc(
			"gostrategy_user_wins",
			"Wins per player in top leaderboard",
			[]string{labelUser}, nil,
		),
		userLossesDesc: prometheus.NewDesc(
			"gostrategy_user_losses",
			"Losses per player in top leaderboard",
			[]string{labelUser}, nil,
		),
		userDrawsDesc: prometheus.NewDesc(
			"gostrategy_user_draws",
			"Draws per player in top leaderboard",
			[]string{labelUser}, nil,
		),
		userGamesDesc: prometheus.NewDesc(
			"gostrategy_user_games",
			"Total games per player in top leaderboard",
			[]string{labelUser}, nil,
		),
		userAvgDurationDesc: prometheus.NewDesc(
			"gostrategy_user_avg_duration_minutes",
			"Average game duration in minutes per player in top leaderboard",
			[]string{labelUser}, nil,
		),
	}
}

// Describe writes all descriptors to the channel.
func (c *DatabaseCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.usersTotalDesc
	ch <- c.usersNew7dDesc
	ch <- c.gamesTotalDesc
	ch <- c.gamesActiveDesc
	ch <- c.gameDurationAvgMinsDesc
	ch <- c.gamesByTypeDesc
	ch <- c.gamesByStatusDesc
	ch <- c.movesTotalDesc
	ch <- c.userWinsDesc
	ch <- c.userLossesDesc
	ch <- c.userDrawsDesc
	ch <- c.userGamesDesc
	ch <- c.userAvgDurationDesc
}

// Collect returns cached metrics or triggers a cache refresh if TTL has expired.
func (c *DatabaseCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	if time.Since(c.lastUpdated) >= c.cacheTTL {
		c.refreshCache()
	}
	metrics := c.cache
	c.mu.Unlock()

	ch <- prometheus.MustNewConstMetric(c.usersTotalDesc, prometheus.GaugeValue, metrics.usersTotal)
	ch <- prometheus.MustNewConstMetric(c.usersNew7dDesc, prometheus.GaugeValue, metrics.usersNew7d)
	ch <- prometheus.MustNewConstMetric(c.gamesTotalDesc, prometheus.GaugeValue, metrics.gamesTotal)
	ch <- prometheus.MustNewConstMetric(c.gamesActiveDesc, prometheus.GaugeValue, metrics.gamesActive)
	ch <- prometheus.MustNewConstMetric(c.gameDurationAvgMinsDesc, prometheus.GaugeValue, metrics.gameDurationAvgMins)
	ch <- prometheus.MustNewConstMetric(c.movesTotalDesc, prometheus.GaugeValue, metrics.movesTotal)

	for gType, val := range metrics.gamesByType {
		ch <- prometheus.MustNewConstMetric(c.gamesByTypeDesc, prometheus.GaugeValue, val, gType)
	}
	for status, val := range metrics.gamesByStatus {
		ch <- prometheus.MustNewConstMetric(c.gamesByStatusDesc, prometheus.GaugeValue, val, status)
	}

	for _, entry := range metrics.leaderboard {
		ch <- prometheus.MustNewConstMetric(c.userWinsDesc, prometheus.GaugeValue, entry.wins, entry.username)
		ch <- prometheus.MustNewConstMetric(c.userLossesDesc, prometheus.GaugeValue, entry.losses, entry.username)
		ch <- prometheus.MustNewConstMetric(c.userDrawsDesc, prometheus.GaugeValue, entry.draws, entry.username)
		ch <- prometheus.MustNewConstMetric(c.userGamesDesc, prometheus.GaugeValue, entry.games, entry.username)
		ch <- prometheus.MustNewConstMetric(c.userAvgDurationDesc, prometheus.GaugeValue, entry.avgDurationMinutes, entry.username)
	}
}

// refreshCache queries the database and populates the local cache structure.
func (c *DatabaseCollector) refreshCache() {
	if db.DB == nil {
		return
	}

	var usersTotal int64
	db.DB.Model(&models.User{}).Where("deleted_at IS NULL").Count(&usersTotal)
	c.cache.usersTotal = float64(usersTotal)

	var usersNew7d int64
	db.DB.Model(&models.User{}).Where("deleted_at IS NULL AND created_at >= ?", time.Now().AddDate(0, 0, -7)).Count(&usersNew7d)
	c.cache.usersNew7d = float64(usersNew7d)

	c.refreshGames()

	var movesTotal int64
	db.DB.Model(&models.GameMove{}).Where("deleted_at IS NULL").Count(&movesTotal)
	c.cache.movesTotal = float64(movesTotal)

	c.refreshLeaderboard()

	c.lastUpdated = time.Now()
}

func (c *DatabaseCollector) refreshGames() {
	var gamesTotal int64
	db.DB.Model(&models.Game{}).Where("deleted_at IS NULL").Count(&gamesTotal)
	c.cache.gamesTotal = float64(gamesTotal)

	gamesActive := getActiveSessionCount()
	c.cache.gamesActive = gamesActive

	var finishedGamesCount int64
	db.DB.Model(&models.Game{}).Where("finished_at IS NOT NULL AND deleted_at IS NULL").Count(&finishedGamesCount)

	c.cache.gamesByStatus = map[string]float64{
		"active":   gamesActive,
		"finished": float64(finishedGamesCount),
	}

	var avgDuration float64
	if finishedGamesCount > 0 {
		if db.DB.Name() == "postgres" {
			_ = db.DB.Model(&models.Game{}).
				Where("finished_at IS NOT NULL AND deleted_at IS NULL").
				Select("ROUND(AVG(EXTRACT(EPOCH FROM (finished_at - created_at))/60)::numeric, 1)").
				Row().Scan(&avgDuration)
		} else {
			_ = db.DB.Model(&models.Game{}).
				Where("finished_at IS NOT NULL AND deleted_at IS NULL").
				Select("AVG(julianday(finished_at) - julianday(created_at)) * 1440").
				Row().Scan(&avgDuration)
		}
	}
	c.cache.gameDurationAvgMins = avgDuration

	// Game Type Split
	type gameTypeSplit struct {
		GameType string
		Count    int64
	}
	var splits []gameTypeSplit
	db.DB.Model(&models.Game{}).
		Where("deleted_at IS NULL").
		Select("game_type, count(*) as count").
		Group("game_type").
		Scan(&splits)

	c.cache.gamesByType = make(map[string]float64)
	for _, split := range splits {
		c.cache.gamesByType[split.GameType] = float64(split.Count)
	}
}

func (c *DatabaseCollector) refreshLeaderboard() {
	type leaderboardDBEntry struct {
		Username               string
		Wins                   int
		Losses                 int
		Draws                  int
		TotalGames             int
		AvgGameDurationSeconds float64
	}
	var dbEntries []leaderboardDBEntry
	db.DB.Table("user_stats").
		Select("users.username, user_stats.wins, user_stats.losses, user_stats.draws, user_stats.total_games, user_stats.avg_game_duration_seconds").
		Joins("JOIN users ON users.id = user_stats.user_id").
		Where("users.deleted_at IS NULL").
		Order("user_stats.wins DESC").
		Limit(10).
		Scan(&dbEntries)

	c.cache.leaderboard = make([]leaderboardEntry, len(dbEntries))
	for i, entry := range dbEntries {
		c.cache.leaderboard[i] = leaderboardEntry{
			username:           entry.Username,
			wins:               float64(entry.Wins),
			losses:             float64(entry.Losses),
			draws:              float64(entry.Draws),
			games:              float64(entry.TotalGames),
			avgDurationMinutes: entry.AvgGameDurationSeconds / 60.0,
		}
	}
}
