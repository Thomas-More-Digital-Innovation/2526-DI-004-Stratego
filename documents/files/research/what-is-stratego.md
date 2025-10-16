# what is stratego?

## objective
capture the **enemy flag** or make the opponent **unable to move**
#### you win if:
1. you capture the flag
2. all their movable pieces are gone

## board
- 10x10 grid
- each player controls one side (4 rows)
- middle has 2 lakes (2x2, non valid positions)
    
## pieces
each player has 40 pieces, hidden from the opponent until battle

| rank | piece name  | count | can move | notes                                                        |
|------|-------------|-------|----------|--------------------------------------------------------------|
| 0    | Flag        | 1     | no       | captures = instant win                                       |
| B    | Bomb        | 6     | no       | destroys any attacker except miner                           |
| 1    | Spy         | 1     | yes      | only piece that can defeat Marshal (10)(if it attacks first) |
| 2    | Scout       | 8     | yes      | can move unlimited in straight line                          |
| 3    | Miner       | 5     | yes      | can defuse a bomb                                            |
| 4    | Sergeant    | 4     | yes      |                                                              |
| 5    | Lieutentant | 4     | yes      |                                                              |
| 6    | Captain     | 4     | yes      |                                                              |
| 7    | Major       | 3     | yes      |                                                              |
| 8    | Colonel     | 2     | yes      |                                                              |
| 9    | General     | 1     | yes      |                                                              |
| 10   | Marshal     | 1     | yes      | The strongest piece                                          |

## combat rules
when 2 pieces meet: 
1. the **attacker** moves onto the defender's square
2. both reveal their ranks
   - Higher rank wins, lower is removed
   - equal rank? -> both are removed
3. special cases
   - Spy defeats marshal when marshal is the attacker
   - miner defuses bomb, other pieces lose against bomb

## movement rules
- most pieces move 1 square at a time, only horizontal & vertical
- scouts can move unlimited in straight line
- you can't move in the lake (not jesus) or on other friendly pieces
- you must move 1 piece each turn

## game loop
1. each player secretly arranges their 40 pieces
2. players take turns to move their pieces
3. battles occur on contact with enemy pieces
4. game ends when:
    - flag is captured
    - player has no legal moves left (or no movable pieces left)

## board

| Row | 0  | 1  | 2  | 3  | 4  | 5  | 6  | 7  | 8  | 9  |
|------|----|----|----|----|----|----|----|----|----|----|
| 0 🔵 | 10 | 9  | 8  | 7  | 6  | 6  | 7  | 8  | 9  | F |
| 1 🔵 | 5  | 5  | 4  | 4  | 3  | 3  | 4  | 4  | 5  | 5  |
| 2 🔵 | 2  | 2  | 2  | 3  | B  | B  | 3  | 2  | 2  | 2  |
| 3 🔵 | S  | 1  | B  | 10  | B  | B  | 3  | 2  | 2  | 2  |
| 4 🟩 |    |    | 🌊 | 🌊 |    |    | 🌊 | 🌊 |    |    |
| 5 🟩 |    |    | 🌊 | 🌊 |    |    | 🌊 | 🌊 |    |    |
| 6 🟥 | S  | 1  | B  | 7  | B  | B  | 3  | 2  | 2  | 2  |
| 7 🟥 | 2  | 2  | 2  | 3  | B  | B  | 3  | 2  | 2  | 2  |
| 8 🟥 | 5  | 5  | 4  | 4  | 3  | 3  | 4  | 4  | 5  | 5  |
| 9 🟥 | 10 | 9  | 8  | F  | 6  | 6  | 7  | 8  | 9  | 10 |

