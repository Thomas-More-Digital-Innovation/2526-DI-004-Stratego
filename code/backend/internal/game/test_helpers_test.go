package game

// setupTestGame initializes the game with players piet and bob and human controllers
func setupTestGame() (*Game, *Player, *Player) {
	player1 := NewPlayer(0, "piet", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(1, "bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	g := NewGame(controller1, controller2)
	return g, &player1, &player2
}

// setupTestSession initializes a session with player controllers piet and bob
func setupTestSession(id string, opts ...SessionOption) (*Session, *Player, *Player) {
	player1 := NewPlayer(0, "piet", "red")
	controller1 := NewHumanPlayerController(&player1)
	player2 := NewPlayer(1, "bob", "blue")
	controller2 := NewHumanPlayerController(&player2)
	session := NewSession(id, controller1, controller2, opts...)
	return session, &player1, &player2
}
