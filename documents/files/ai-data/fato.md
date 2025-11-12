# FATO AI

FATO AI is a random-move agent based AI who can remember the board & act on it. It does not have the ability to learn anything as it is still random move based. It inherits the random movement from Fafo AI with new logic that checks the board for known enemy pieces and acts on that knowledge.

> FATO means "Fuck Around & Try Out"

## Results

We put 2 FATO AI players against each other.
Alice and Bob, both take equal number of turns and use random board setups.

> go run . --ai=fato:fato --format md --logging=false --matches {n}

### AI vs AI Tournament Summary (10 games)

**Total Matches:** 10
**Total Rounds:** 2996
**Average Rounds (per game):** 299.60
**Shortest Game (rounds):** 34

#### Overall Win Causes

|             Cause | Count |     % |
| ----------------: | ----: | ----: |
|     Flag captured |     4 | 80.0% |
| No movable pieces |     1 | 20.0% |
|         Max turns |     0 |  0.0% |

#### Player Results

| Player          | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
| :-------------- | ---: | ----: | ------------: | -----------: | ------------: |
| Alice AI - fato |    2 | 40.0% |             2 |            0 |             0 |
| Bob AI - fato   |    3 | 60.0% |             2 |            1 |             0 |

**Draws:** 5 (100.0%)

AI vs AI matches completed in 0.06 seconds

### AI vs AI Tournament Summary (100 games)

**Total Matches:** 100
**Total Rounds:** 22839
**Average Rounds (per game):** 228.39
**Shortest Game (rounds):** 1

#### Overall Win Causes

|             Cause | Count |     % |
| ----------------: | ----: | ----: |
|     Flag captured |    61 | 77.2% |
| No movable pieces |     8 | 10.1% |
|         Max turns |    10 | 12.7% |

#### Player Results

| Player          | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
| :-------------- | ---: | ----: | ------------: | -----------: | ------------: |
| Alice AI - fato |   38 | 48.1% |            29 |            5 |             4 |
| Bob AI - fato   |   41 | 51.9% |            32 |            3 |             6 |

**Draws:** 21 (26.6%)

AI vs AI matches completed in 0.55 seconds

### AI vs AI Tournament Summary (1000 games)

**Total Matches:** 1000  
**Total Rounds:** 244144  
**Average Rounds (per game):** 244.14  
**Shortest Game (rounds):** 1

#### Overall Win Causes

| Cause | Count | % |
|-------:|------:|---:|
| Flag captured | 502 | 67.2% |
| No movable pieces | 138 | 18.5% |
| Max turns | 107 | 14.3% |

#### Player Results

| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
|:-------|-----:|-----:|--------------:|-------------:|--------------:|
| Alice AI - fato | 372 | 49.8% | 253 | 67 | 52 |
| Bob AI - fato | 375 | 50.2% | 249 | 71 | 55 |

**Draws:** 253 (33.9%)

### AI vs AI Tournament Summary (10000 games)

**Total Matches:** 10000  
**Total Rounds:** 2381773  
**Average Rounds (per game):** 238.18  
**Shortest Game (rounds):** 1

#### Overall Win Causes

| Cause | Count | % |
|-------:|------:|---:|
| Flag captured | 5328 | 70.8% |
| No movable pieces | 1090 | 14.5% |
| Max turns | 1111 | 14.8% |

#### Player Results

| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
|:-------|-----:|-----:|--------------:|-------------:|--------------:|
| Alice AI - fato | 3691 | 49.0% | 2610 | 550 | 531 |
| Bob AI - fato | 3838 | 51.0% | 2718 | 540 | 580 |

**Draws:** 2471 (32.8%)

AI vs AI matches completed in 57.78 seconds
