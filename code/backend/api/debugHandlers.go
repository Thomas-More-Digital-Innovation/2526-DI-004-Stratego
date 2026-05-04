package api

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

// DebugStats returns runtime memory and goroutine statistics
func (s *GameServer) DebugStats(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s.mutex.RLock()
	sessionCount := len(s.sessions)
	s.mutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"goroutines":    runtime.NumGoroutine(),
		"activeSessions": sessionCount,
		"memory": gin.H{
			"alloc":      m.Alloc / 1024 / 1024,
			"totalAlloc": m.TotalAlloc / 1024 / 1024,
			"sys":        m.Sys / 1024 / 1024,
			"numGC":      m.NumGC,
		},
	})
}
