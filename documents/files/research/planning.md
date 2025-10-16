# Phase 1 — Research & Design
- Study official Stratego rules and edge cases (combat, movement, victory).
- Define exact data structures: board size, coordinates, piece types, ranks.

Plan internal modules:
- engine → core game logic 
- rules → game rules
- ai → computer player(s)
- sim → self-play and testing
- server → later for sockets/front-end

Deliverable → a written spec and directory layout.

# Phase 2 — Core Game Engine
- Implement all rules in isolation: setup, legal move generation, battle resolution, win detection.
- Include functions to clone, apply, and undo moves (needed for AI search).
- Add deterministic random seed handling for reproducible simulations.
- Write unit tests for every rule and scenario.

Deliverable → fully working headless game engine tested via CLI or automated tests.

# Phase 3 — AI Foundation
- Start with a baseline random-move agent for validation.
- Design a generic Agent interface so you can swap AI types easily.
- Add simple heuristic evaluation (piece value, mobility, safety).

Deliverable → two random AIs can play a full legal match.

# Phase 4 — Search-Based AI
- Implement minimax with alpha-beta pruning using your clone/undo functions.
- Add move ordering and adjustable depth.
- Compare heuristics and tune evaluation weights.
- Log thousands of AI-vs-AI games to gather performance data.

Deliverable → a competitive AI capable of consistent logical play.

# Phase 5 — Advanced AI Exploration
- Experiment with Monte Carlo Tree Search or hybrid minimax + probability approaches.
- Optionally add learning through self-play data (reinforcement-like).
- Optimize with concurrency and caching (transposition tables).

Deliverable → “strong” AI engine ready for human challenge.

# Phase 6 — Networking & Server Layer
- Build a WebSocket-based game server exposing create/join/move/update events.
- Maintain active sessions and synchronize turns.
- Keep separate REST endpoints for setup and stats.

Deliverable → headless multiplayer backend usable by GUI or AI clients.

# Phase 7 — Frontend (Optional Later)
- Create a lightweight Svelte interface showing the board and moves.
- Connect via WebSocket to play vs AI or another player.

Deliverable → playable browser version backed by your Go engine.

# Phase 8 — Testing & Benchmarking
- Stress-test thousands of AI matches to measure win rates and speed.
- Profile CPU/memory to refine performance.

Ensure deterministic replays for debugging.

# Phase 9 — Release & Iteration
- Package engine as a library or standalone binary.
- Open API for community AI plugins.

Publish docs and examples.
