# Different AIs for Stratego

## 🎲 1. **FAFO (Random)**

**Idea:**
Moves are chosen randomly from all legal actions.

**✅ Pros:**

* Easiest to implement.
* Provides a good baseline for testing.
* Always unpredictable.

**❌ Cons:**

* No learning or strategy.
* Performs very poorly against anything smarter.

**🧠 Trainable?**
Not really — random agents don’t learn. They can, however, be used as “opponents” for training smarter agents (e.g., reinforcement learning warm-up).

---

## ⚙️ 2. **Heuristic**

**Idea:**
Uses manually designed rules or evaluations (e.g., value pieces, favor capturing lower-ranked enemies, protect flag).

**✅ Pros:**

* Fast and simple.
* Encodes expert knowledge easily.
* Can perform decently without heavy computation.

**❌ Cons:**

* Limited by human bias or oversimplified rules.
* Hard to adapt to hidden information and deception in Stratego.
* Doesn’t improve without manual tuning.

**🧠 Trainable?**
Semi-trainable — you can optimize weights of heuristics (e.g., via genetic algorithms or reinforcement learning) to improve over time.

---

## ♟️ 3. **Minimax**

**Idea:**
Explores the game tree assuming both players play optimally. Each node alternates between maximizing and minimizing the evaluation score.

**✅ Pros:**

* Theoretically strong — finds optimal play if the tree is fully explored.
* Great for deterministic perfect-information games (like chess).

**❌ Cons:**

* Stratego has **hidden information** (unknown opponent pieces), so minimax can’t model uncertainty well.
* The branching factor is huge → needs pruning (α–β pruning) and depth limits.
* Struggles when bluffing or incomplete knowledge is key.

**🧠 Trainable?**
Partly — you can train the **evaluation function** (e.g., using self-play to learn board value estimates). But minimax’s structure itself is not learnable.

---

## 🌳 4. **MCTS (Monte Carlo Tree Search)**

**Idea:**
Simulates many random playouts from the current position to estimate move quality statistically. Expands the tree towards promising moves using exploration/exploitation balance.

**✅ Pros:**

* Handles huge and uncertain state spaces better than minimax.
* Adapts dynamically — no fixed evaluation needed.
* Excellent for **hidden-information games** (if you include belief modeling).
* Basis of AlphaZero-style learning.

**❌ Cons:**

* Computationally heavy (many simulations).
* Quality depends on playout policy (random = weak, learned = stronger).
* Requires many iterations for stable results.

**🧠 Trainable?**
Yes — very trainable. You can:

* Train a **policy network** to guide simulations.
* Train a **value network** to replace random rollouts.
* Use **self-play reinforcement learning** (AlphaZero-style) to improve both.

---

## ⚔️ TL;DR — in Stratego context

| Algorithm | Info type     | Strength         | Weakness             | Trainable  | Notes                 |
| --------- | ------------- | ---------------- | -------------------- | ---------- | --------------------- |
| Random    | None          | Unpredictable    | Dumb                 | ❌          | Use for testing       |
| Heuristic | Expert rules  | Fast             | Rigid, biased        | ⚙️ Partial | Good baseline         |
| Minimax   | Deterministic | Strategic        | Hidden info kills it | ⚙️ Partial | Needs belief modeling |
| MCTS      | Statistical   | Flexible, strong | Heavy compute        | ✅ Full     | Best long-term option |

## Expected runtime difference between FAFO & MCTS

* **Baseline:** 1.44 s for **1000 games total** (so **0.00144 s per game**).
* **moves/game** = total moves (both players combined). I show 50, 100, 200.
* Random agent = 1 “eval” per move; MCTS uses **S** simulations per decision.
* One MCTS vs random → MCTS player makes about half the decisions.

### Key results (rounded)

For **S = 100 sims / decision**

* **One MCTS vs random:** **72.72 s** total for 1000 games ≈ **1.21 min** (≈ **50.5×** slower than baseline).
* **Both MCTS:** **144.0 s** ≈ **2.40 min** (≈ **100×** slower).

For **S = 1,000**

* **One MCTS:** **720.72 s** ≈ **12.01 min** (≈ **500.5×**).
* **Both MCTS:** **1,440 s** = **24.00 min** (≈ **1000×**).

For **S = 10,000**

* **One MCTS:** **7,200.72 s** ≈ **2.00 hours** (≈ **5000.5×**).
* **Both MCTS:** **14,400 s** = **4.00 hours** (≈ **10000×**).

**TL;DR:** Even with the tiny baseline (1.44 s / 1000 games), MCTS blows up quickly. Cost scales roughly **linearly** with simulations (S) and decisions per game.
**If runtime matters:** lower S, parallelize sims, use a learned policy/value to cut sims, or only use MCTS on high-value turns.
