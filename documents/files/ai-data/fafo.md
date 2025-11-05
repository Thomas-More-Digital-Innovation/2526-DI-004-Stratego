# FAFO AI

FAFO AI is a simple random-move agent.

> FAFO means "Fuck Around & Find Out"

## Results

We put 2 FAFO AI players against each other.
Alice and Bob, both take equal number of turns and use random board setups.

> go run . --format md --logging=false --matches {n}

### 📊 Tournament Results (10 games)

| Player   | Wins | Win % |
| -------- | ---- | ----- |
| Alice AI | 5    | 50.0% |
| Bob AI   | 4    | 40.0% |
| Draws    | 1    | 10.0% |

| Win Causes        | Count | %     |
| ----------------- | ----- | ----- |
| Flag captured     | 8     | 80.0% |
| No movable pieces | 0     | 0.0%  |
| Max turns         | 2     | 20.0% |

Average game length: 169.8 rounds
AI vs AI matches completed in 0.01 seconds

### 📊 Tournament Results (100 games)

| Player   | Wins | Win % |
| -------- | ---- | ----- |
| Alice AI | 42   | 42.0% |
| Bob AI   | 45   | 45.0% |
| Draws    | 13   | 13.0% |

| Win Causes        | Count | %     |
| ----------------- | ----- | ----- |
| Flag captured     | 77    | 77.0% |
| No movable pieces | 1     | 1.0%  |
| Max turns         | 22    | 22.0% |

Average game length: 184.6 rounds

AI vs AI matches completed in 0.13 seconds

### 📊 Tournament Results (1000 games)

| Player   | Wins | Win % |
| -------- | ---- | ----- |
| Alice AI | 421  | 42.1% |
| Bob AI   | 458  | 45.8% |
| Draws    | 121  | 12.1% |

| Win Causes        | Count | %     |
| ----------------- | ----- | ----- |
| Flag captured     | 770   | 77.0% |
| No movable pieces | 6     | 0.6%  |
| Max turns         | 224   | 22.4% |

Average game length: 172.3 rounds

AI vs AI matches completed in 1.44 seconds

### 📊 Tournament Results (10000 games)

| Player   | Wins | Win % |
| -------- | ---- | ----- |
| Alice AI | 4346 | 43.5% |
| Bob AI   | 4309 | 43.1% |
| Draws    | 1345 | 13.4% |

| Win Causes        | Count | %     |
| ----------------- | ----- | ----- |
| Flag captured     | 7383  | 73.8% |
| No movable pieces | 109   | 1.1%  |
| Max turns         | 2508  | 25.1% |

Average game length: 183.5 rounds

AI vs AI matches completed in 15.27 seconds

### 📊 Tournament Results (100000 games)

| Player   | Wins  | Win % |
| -------- | ----- | ----- |
| Alice AI | 43020 | 43.0% |
| Bob AI   | 43577 | 43.6% |
| Draws    | 13403 | 13.4% |

| Win Causes        | Count | %     |
| ----------------- | ----- | ----- |
| Flag captured     | 74013 | 74.0% |
| No movable pieces | 1080  | 1.1%  |
| Max turns         | 24907 | 24.9% |

Average game length: 183.1 rounds

AI vs AI matches completed in 334.64 seconds

## Analysis

Over increasing match counts, both **Alice AI** and **Bob AI** show nearly equal performance, confirming that the AIs’ random strategies and setups produce balanced outcomes.

### 🏆 Performance Trends

* **Small samples (10–100 games)** show slight random variation — Alice led early on, but as matches increased, Bob gained a small edge.
* **At 100,000 games**, Bob leads marginally with **43.6% wins** versus Alice’s **43.0%**, while **draws stabilize around 13–14%**.
* These results suggest **statistical parity** — no significant advantage for either AI.

### 🎯 Win Conditions

* The **most common victory** is **flag capture (~74–77%)**, showing that most games end through standard play rather than exhaustion.
* **Stalemates (max turns reached)** grow slightly with more games, from 20% to ~25%, indicating long, indecisive matches are common in random setups.
* **No movable pieces** remains rare (~1%), suggesting both AIs rarely trap themselves completely.

### ⏱️ Game Duration

* Average length stabilizes around **180 rounds**, meaning random strategies lead to lengthy battles.
* Despite the size of the tournaments, **simulation speed scales efficiently** — 100,000 matches finish in ~335 seconds, showing strong performance of the game engine.

### 📈 Summary

* Both AIs perform **evenly** across all scales.
* **Random play** leads to balanced results but **long and often unresolved games**.
* The **flag capture** remains the dominant end condition, reinforcing that despite randomness, most games reach a natural conclusion.

---

## Summary

Alice and Bob are equally matched random players. Most wins come from flag captures, with ~13% draws and long games averaging 180 rounds.
