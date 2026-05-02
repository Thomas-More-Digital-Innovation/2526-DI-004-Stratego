# Changelog

All notable changes to this project will be documented in this file.

## [0.2.0] - 02-05-2026

### Added
- **Database Migrations:** Implemented a robust SQL migration system using Go's `embed` filesystem.
- **Session Persistence:** Games and moves are now automatically saved to the database.
- **User Stats:** Automated tracking of wins, losses, and move counts for registered users.
- **Setup Phase Timeout:** Added a 5-minute automated setup phase with warnings to prevent stalled games.
- **JWT v5 Migration:** Upgraded security infrastructure to use the latest `golang-jwt/v5` standard.
- **IP Rate Limiting:** Added protection against brute-force attacks on authentication endpoints.

### Changed
- **Concurrency Overhaul:** Optimized `GameRunner` and `WSHub` with event-driven channels, eliminating busy-waiting and reducing mutex contention.
- **Standardized Logging:** Refactored logging to use a consistent `[Time] [Tag] [Location] [User] Message` format.
- **AI Turn Pacing:** Improved AI move synchronization to ensure smooth animations in the frontend.

### Fixed
- Illegal piece movement validation (e.g., Flags can no longer move).
- Game state broadcasting regressions during high-concurrency sessions.
- Authorization leaks where spectators could occasionally attempt to send moves.

## [0.1.2] - 13-04-2026

### Fixed
- Mutex contention issues identified in system audit.
- Empty CORS origin handling.
- API path parameter inconsistencies.

## [0.1.1] - 03-04-2026

### Added
- Initial functional release of Stratego with Human-vs-AI and AI-vs-AI modes.
- Basic WebSocket implementation for game updates.
- In-memory session management.
