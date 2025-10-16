# Core Game & AI Engine
- Go

stdlibs -> net/http, encoding/json, math/rand, sync
unittest -> testing
profiling (optimizing AI) -> pprof

/engine   → rules, board, moves, state
/ai       → random, heuristic, minimax, mcts
/sim      → self-play and evaluation
/server   → HTTP + WebSocket handlers

# network layer
make go backend communicate with browsers & ai clients

### protocols
- websockets -> live play, move updates, chat?
- HTTP/REST -> setup, listing games, stats, ...

libs: 
- HTTP router -> chi
- websocket -> nhooyr.io/websocket
- logs -> rs/zerolog

# frontend
Svelte + Typescript (NO JAVASCRIPT!!!)

PixiJS or Phaser for 2D graphics?
