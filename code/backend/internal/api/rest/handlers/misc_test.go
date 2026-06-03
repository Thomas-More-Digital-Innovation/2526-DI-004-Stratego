package handlers_test

import (
	"context"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/pkg/testutils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/health", nil)
	h.HealthHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCSRFToken(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h.GetCSRFToken(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStatsHandlers(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

	user, err := db.CreateUser(context.Background(), "alice", "StrongPassword1", "")
	require.NoError(t, err)

	// GetUserStatsHandler
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("GET", fmt.Sprintf("/stats/%d", user.ID), nil)
	c1.Params = []gin.Param{{Key: "id", Value: fmt.Sprintf("%d", user.ID)}}
	h.GetUserStatsHandler(c1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// GetCurrentUserStatsHandler
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("GET", "/stats/me", nil)
	c2.Set("user", user)
	h.GetCurrentUserStatsHandler(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestDebugStats(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h.DebugStats(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMonitoringHandlers(t *testing.T) {
	_, h, _ := testutils.SetupHandlerTest(t)

	// UserCountHandler
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request, _ = http.NewRequest("GET", "/monitoring/users/count", nil)
	h.UserCountHandler(c1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// GamesPlayedCountHandler
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request, _ = http.NewRequest("GET", "/monitoring/games/count", nil)
	h.GamesPlayedCountHandler(c2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
