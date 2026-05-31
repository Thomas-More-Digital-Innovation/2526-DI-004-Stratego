package telemetry

import "sync/atomic"

var activeSessionCount int64

// TrackSession increments the active game session counter atomically.
func TrackSession() { atomic.AddInt64(&activeSessionCount, 1) }

// UntrackSession decrements the active game session counter atomically.
func UntrackSession() { atomic.AddInt64(&activeSessionCount, -1) }

func getActiveSessionCount() float64 { return float64(atomic.LoadInt64(&activeSessionCount)) }
