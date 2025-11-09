# FATO AI

FATO AI is a random-move agent based AI who can remember the board & act on it. It does not have the ability to learn anything as it is still random move based. It inherits the random movement from Fafo AI with new logic that checks the board for known enemy pieces and moves toward them if found.

> FATO means "Fuck Around & Try Out"

## Results

We put 2 FATO AI players against each other.
Alice and Bob, both take equal number of turns and use random board setups.

> go run . --ai=fato:fato --format md --logging=false --matches {n}

### AI vs AI Tournament Summary (10 games)

**Total Matches:** 10  
**Total Rounds:** 1523  
**Average Rounds (per game):** 152.30  
**Shortest Game (rounds):** 1

#### Overall Win Causes

| Cause | Count | % |
|-------:|------:|---:|
| Flag captured | 8 | 88.9% |
| No movable pieces | 0 | 0.0% |
| Max turns | 1 | 11.1% |

#### Player Results

| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
|:-------|-----:|-----:|--------------:|-------------:|--------------:|
| Alice AI - fato | 4 | 44.4% | 4 | 0 | 0 |
| Bob AI - fato | 5 | 55.6% | 4 | 0 | 1 |

**Draws:** 1 (11.1%)

AI vs AI matches completed in 0.02 seconds

### AI vs AI Tournament Summary (100 games)

**Total Matches:** 100  
**Total Rounds:** 18604  
**Average Rounds (per game):** 186.04  
**Shortest Game (rounds):** 1

#### Overall Win Causes

| Cause | Count | % |
|-------:|------:|---:|
| Flag captured | 74 | 85.1% |
| No movable pieces | 0 | 0.0% |
| Max turns | 13 | 14.9% |

#### Player Results

| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
|:-------|-----:|-----:|--------------:|-------------:|--------------:|
| Alice AI - fato | 50 | 57.5% | 41 | 0 | 9 |
| Bob AI - fato | 37 | 42.5% | 33 | 0 | 4 |

**Draws:** 13 (14.9%)

AI vs AI matches completed in 0.22 seconds

### AI vs AI Tournament Summary (1000 games)

**Total Matches:** 1000  
**Total Rounds:** 181686  
**Average Rounds (per game):** 181.69  
**Shortest Game (rounds):** 1

#### Overall Win Causes

| Cause | Count | % |
|-------:|------:|---:|
| Flag captured | 753 | 85.6% |
| No movable pieces | 9 | 1.0% |
| Max turns | 118 | 13.4% |

#### Player Results

| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
|:-------|-----:|-----:|--------------:|-------------:|--------------:|
| Alice AI - fato | 425 | 48.3% | 362 | 2 | 61 |
| Bob AI - fato | 455 | 51.7% | 391 | 7 | 57 |

**Draws:** 120 (13.6%)

AI vs AI matches completed in 2.95 seconds

### AI vs AI Tournament Summary (10000 games)

**Total Matches:** 10000  
**Total Rounds:** 1823626  
**Average Rounds (per game):** 182.36  
**Shortest Game (rounds):** 1

#### Overall Win Causes

| Cause | Count | % |
|-------:|------:|---:|
| Flag captured | 7394 | 85.1% |
| No movable pieces | 106 | 1.2% |
| Max turns | 1191 | 13.7% |

#### Player Results

| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
|:-------|-----:|-----:|--------------:|-------------:|--------------:|
| Alice AI - fato | 4386 | 50.5% | 3738 | 54 | 594 |
| Bob AI - fato | 4305 | 49.5% | 3656 | 52 | 597 |

**Draws:** 1309 (15.1%)

AI vs AI matches completed in 38.82 seconds

### AI vs AI Tournament Summary (100000 games)

**Total Matches:** 100000  
**Total Rounds:** 18250244  
**Average Rounds (per game):** 182.50  
**Shortest Game (rounds):** 1

#### Overall Win Causes

| Cause | Count | % |
|-------:|------:|---:|
| Flag captured | 74139 | 85.4% |
| No movable pieces | 1105 | 1.3% |
| Max turns | 11527 | 13.3% |

#### Player Results

| Player | Wins | Win % | Flag captures | No-move wins | Max-turn wins |
|:-------|-----:|-----:|--------------:|-------------:|--------------:|
| Alice AI - fato | 43264 | 49.9% | 36883 | 533 | 5848 |
| Bob AI - fato | 43507 | 50.1% | 37256 | 572 | 5679 |

**Draws:** 13229 (15.2%)

AI vs AI matches completed in 331.41 seconds
