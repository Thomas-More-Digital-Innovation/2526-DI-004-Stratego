package engine_test

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"testing"
)

func TestOwner(t *testing.T) {
	player1 := engine.NewPlayer(1, "player1", "avatar1")
	player2 := engine.NewPlayer(2, "player2", "avatar2")
	piece1 := engine.NewPiece(models.Flag, &player1)
	piece2 := engine.NewPiece(models.Flag, &player2)

	if *piece1.GetOwner() != player1 {
		t.Errorf("P1 Expected owner to be %v, got %v", player1, *piece1.GetOwner())
	}
	if *piece2.GetOwner() != player2 {
		t.Errorf("P2 Expected owner to be %v, got %v", player2, *piece2.GetOwner())
	}

	if piece1.GetOwner().GetID() == piece2.GetOwner().GetID() {
		t.Errorf("Expected different owners, got %v and %v", piece1.GetOwner(), piece2.GetOwner())
	}
}

func TestType(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	piece1 := engine.NewPiece(models.Flag, &player)
	piece2 := engine.NewPiece(models.Bomb, &player)

	if piece1.GetType() == piece2.GetType() {
		t.Errorf("Expected different piece types, got %v and %v", piece1.GetType(), piece2.GetType())
	}
}

func TestGetRank(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	flag := engine.NewPiece(models.Flag, &player)
	bomb := engine.NewPiece(models.Bomb, &player)
	scout := engine.NewPiece(models.Scout, &player)
	marshal := engine.NewPiece(models.Marshal, &player)

	if flag.GetRank() != '0' {
		t.Errorf("Expected rank to be '0', got %d", flag.GetRank())
	}

	if bomb.GetRank() != 'B' {
		t.Errorf("Expected rank to be 'B', got %d", bomb.GetRank())
	}

	if scout.GetRank() != '2' {
		t.Errorf("Expected rank to be '2', got %d", scout.GetRank())
	}

	if marshal.GetRank() != 'M' {
		t.Errorf("Expected rank to be 'M', got %d", marshal.GetRank())
	}
}

func TestCanMove(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	flag := engine.NewPiece(models.Flag, &player)
	bomb := engine.NewPiece(models.Bomb, &player)
	general := engine.NewPiece(models.General, &player)

	if flag.CanMove() {
		t.Errorf("Expected flag to not be able to move, got true")
	}
	if bomb.CanMove() {
		t.Errorf("Expected bomb to not be able to move, got true")
	}
	if !general.CanMove() {
		t.Errorf("Expected general to be able to move, got false")
	}
}

func TestRevealed(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	spy := engine.NewPiece(models.Spy, &player)

	if spy.IsRevealed() {
		t.Errorf("Expected spy to not be revealed, got true")
	}

	spy.Reveal()

	if !spy.IsRevealed() {
		t.Errorf("Expected spy to be revealed, got false")
	}
}

func TestEliminated(t *testing.T) {
	player := engine.NewPlayer(1, "player1", "avatar1")
	major := engine.NewPiece(models.Major, &player)
	player.InitializePieceScore(major.GetStrategicValue())

	if !major.IsAlive() {
		t.Errorf("Expected major to be alive, got false")
	}

	major.Eliminate()

	if major.IsAlive() {
		t.Errorf("Expected major to be eliminated, got true")
	}

	expectedScore := 0
	if player.GetPieceScore() != expectedScore {
		t.Errorf("Expected player piece score to be %d, got %d", expectedScore, player.GetPieceScore())
	}
}

func TestSpyAttackingMarshal(t *testing.T) {
	// attacker
	player1 := engine.NewPlayer(1, "player1", "avatar1")
	spy := engine.NewPiece(models.Spy, &player1)
	player1.InitializePieceScore(spy.GetStrategicValue())

	// target
	player2 := engine.NewPlayer(2, "player2", "avatar2")
	marshal := engine.NewPiece(models.Marshal, &player2)
	player2.InitializePieceScore(marshal.GetStrategicValue())

	result := spy.Attack(marshal)
	attacker, target := result[0], result[1]

	if target.IsAlive() {
		t.Errorf("Expected marshal to be eliminated, got alive")
	}
	if !attacker.IsAlive() {
		t.Errorf("Expected spy to be alive, got eliminated")
	}

	expectedScore1 := models.Spy.GetStrategicValue()
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player piece (spy) score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := 0
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}
}

func TestMinerAttackingBomb(t *testing.T) {
	// setup
	player1 := engine.NewPlayer(1, "player1", "avatar1")
	miner := engine.NewPiece(models.Miner, &player1)
	player1.InitializePieceScore(miner.GetStrategicValue())

	player2 := engine.NewPlayer(2, "player2", "avatar2")
	bomb := engine.NewPiece(models.Bomb, &player2)
	player2.InitializePieceScore(bomb.GetStrategicValue())

	// execute
	result := miner.Attack(bomb)
	miner, bomb = result[0], result[1]

	// verify
	if bomb.IsAlive() {
		t.Errorf("Expected bomb to be eliminated, got alive")
	}
	if !miner.IsAlive() {
		t.Errorf("Expected miner to be alive, got eliminated")
	}

	expectedScore1 := models.Miner.GetStrategicValue()
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player1 (miner) score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := 0
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 (bomb) piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}
}

func TestAttackBomb(t *testing.T) {
	// setup
	player1 := engine.NewPlayer(1, "player1", "avatar1")
	colonel := engine.NewPiece(models.Colonel, &player1)
	player1.InitializePieceScore(colonel.GetStrategicValue())

	player2 := engine.NewPlayer(2, "player2", "avatar2")
	bomb := engine.NewPiece(models.Bomb, &player2)
	player2.InitializePieceScore(bomb.GetStrategicValue())

	// execute
	result := colonel.Attack(bomb)
	attacker, target := result[0], result[1]

	// verify
	if !target.IsAlive() {
		t.Errorf("Expected bomb to be alive, got eliminated")
	}
	if attacker.IsAlive() {
		t.Errorf("Expected colonel to be eliminated, got alive")
	}
	expectedScore1 := 0
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player1 piece score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := models.Bomb.GetStrategicValue()
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}

}

func TestStandardAttack(t *testing.T) {
	// setup
	player1 := engine.NewPlayer(1, "player1", "avatar1")
	major := engine.NewPiece(models.Major, &player1)
	player1.InitializePieceScore(major.GetStrategicValue())

	player2 := engine.NewPlayer(2, "player2", "avatar2")
	captain := engine.NewPiece(models.Captain, &player2)
	player2.InitializePieceScore(captain.GetStrategicValue())

	// execute
	result := major.Attack(captain)
	attacker, target := result[0], result[1]

	// verify
	if target.IsAlive() {
		t.Errorf("Expected captain to be eliminated, got alive")
	}
	if !attacker.IsAlive() {
		t.Errorf("Expected major to be alive, got eliminated")
	}

	expectedScore1 := models.Major.GetStrategicValue()
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player1 piece (major) score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := 0
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}
}

func TestStandardAttackEqualRank(t *testing.T) {
	// setup
	player1 := engine.NewPlayer(1, "player1", "avatar1")
	sergeant1 := engine.NewPiece(models.Sergeant, &player1)
	player1.InitializePieceScore(sergeant1.GetStrategicValue())

	player2 := engine.NewPlayer(2, "player2", "avatar2")
	sergeant2 := engine.NewPiece(models.Sergeant, &player2)
	player2.InitializePieceScore(sergeant2.GetStrategicValue())

	// execute
	result := sergeant1.Attack(sergeant2)
	attacker, target := result[0], result[1]

	// verify
	if attacker.IsAlive() {
		t.Errorf("Expected attacker sergeant to be eliminated, got alive")
	}
	if target.IsAlive() {
		t.Errorf("Expected target sergeant to be eliminated, got alive")
	}

	expectedScore1 := 0
	if player1.GetPieceScore() != expectedScore1 {
		t.Errorf("Expected player1 piece score to be %d, got %d", expectedScore1, player1.GetPieceScore())
	}

	expectedScore2 := 0
	if player2.GetPieceScore() != expectedScore2 {
		t.Errorf("Expected player2 piece score to be %d, got %d", expectedScore2, player2.GetPieceScore())
	}
}
