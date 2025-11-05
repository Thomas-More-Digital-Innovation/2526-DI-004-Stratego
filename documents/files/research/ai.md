# Different AIs for Stratego

## 🎲 1. **Random**

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
