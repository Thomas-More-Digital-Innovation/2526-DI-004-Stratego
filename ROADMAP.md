# GoStrategy Roadmap

This document outlines the planned improvements and future directions for the GoStrategy project, prioritized by security, infrastructure resilience, and user experience.

---

## v0.4.0 - AI & Intelligence
*Focus: Moving beyond random moves and providing tools for AI experimentation.*

### AI Improvements
- [x] **Heuristics**: Create a basic heuristic-based evaluation engine.
- [x] **Minimax**: Create a minimax-based evaluation engine.
- [x] **MCTS**: Create a Monte Carlo Tree Search-based evaluation engine.
- [ ] **AI In Game**: Add the new AIs to the game as available opponents.
- [ ] **AI Info**: Add an info page about each AIs algorithm, strengths and weaknesses, and how they work.
- [ ] **Export Tools**: Add ability to export board setups as JSON for external AI training/testing.


## v0.4.1 - Replays
*Focus: Storing and viewing past games.*

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
