# GoStrategy Roadmap

This document outlines the planned improvements and future directions for the GoStrategy project, prioritized by security, infrastructure resilience, and user experience.

---

### Current Phase: v0.3.1 - Refactoring & Polish
*Focus: Strengthening the core experience and preparing for advanced features.*

#### Backend & Infrastructure
- [x] **Refactoring**: refactoring of backend code, especially the api package.
- [ ] **Rate Limiting**: Add secondary user-level rate limits for sensitive actions (creation, password changes).
- [x] **Security**: Implement stricter regex-based input sanitization for usernames to prevent injection.
- [ ] **Reliability**: Implement exponential backoff for WebSocket reconnections with UI feedback.
- [x] **Game Limiting**: Allow only one active game per user at a time.

#### Frontend & UI/UX
- [ ] **Authentication**: Redesign and improve the login and registration screens for a more premium feel.
- [ ] **Accessibility (a11y)**: Complete an audit for keyboard navigation and ARIA labeling on the game board.
- [ ] **Feedback**: Add "Reconnecting..." status indicators for socket drops.
- [ ] **Reconnect Button**: As long as a game is waiting for cleanup, allow the user to reconnect to it via a button in the home page.

---

## v0.4.0 - AI & Intelligence
*Focus: Moving beyond random moves and providing tools for AI experimentation.*

### AI Playground Improvements
- [ ] **Dataset Generation**: Finalize automated AI vs AI testing environment for large-scale data collection.
- [ ] **Export Tools**: Add ability to export board setups as JSON for external AI training/testing.
- [ ] **Move History**: Implement a dedicated history view for AI moves in the playground.
- [ ] **Heuristics**: Create and train a basic heuristic-based evaluation engine.

### Replay System
- [ ] **Game Theater**: View replays of previous games with full playback controls.
- [ ] **External Loading**: Support for uploading and viewing replay files.
- [ ] **Lifecycle**: Automated cleanup cycles to manage storage (keep last X games per user).

---

## v0.5.0 - Social & PvP
*Focus: Connecting players and building a competitive community.*

### Social Features
- [ ] **Friendship System**: Friend requests and a dedicated social tab.
- [ ] **Direct Challenges**: Challenge friends to private games.
- [ ] **Matchmaking**: Implementation of a global queue and Elo-based matching.
- [ ] **Social Sharing**: Generate unique game links to invite external players.

---

## v0.6.0 - Mobile & PWA
*Focus: Making GoStrategy available everywhere.*

### Mobile Experience
- [ ] **Responsive Engine**: Full UI overhaul to support touch interfaces and smaller screens.
- [ ] **PWA Implementation**: Service worker support for offline caching and "Add to Home Screen" capability.
- [ ] **Native Feel**: Add PWA splash screens and optimized asset loading.
