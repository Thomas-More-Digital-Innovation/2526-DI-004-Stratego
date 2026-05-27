package models

// Exported piece types for use in the engine and UI
var (
	// Flag is the ultimate target of the game
	Flag = *NewPieceType("Flag", '0', false, "The piece you must capture to win the game.", "🚩", 1, 0)
	// Bomb is a stationary piece that eliminates most attackers
	Bomb = *NewPieceType("Bomb", 'B', false, "The piece that cannot move and eliminates most attackers.", "💣", 6, 7)
	// Spy is the only piece that can defeat the Marshal if it attacks first
	Spy = *NewPieceType("Spy", '1', true, "The piece that can move and attack but is weak.", "🕵️", 1, 7)
	// Scout can move any number of squares in a straight line
	Scout = *NewPieceType("Scout", '2', true, "The piece that can move multiple spaces and attack.", "🕵️‍♂️", 8, 3)
	// Miner can defuse Bombs
	Miner = *NewPieceType("Miner", '3', true, "The piece that can move and attack but is weak.", "⛏️", 5, 6)
	// Sergeant is a rank 4 piece
	Sergeant = *NewPieceType("Sergeant", '4', true, "The piece that can move and attack but is weak.", "👮", 4, 4)
	// Lieutenant is a rank 5 piece
	Lieutenant = *NewPieceType("Lieutenant", '5', true, "The piece that can move and attack but is weak.", "👮‍♂️", 4, 5)
	// Captain is a rank 6 piece
	Captain = *NewPieceType("Captain", '6', true, "The piece that can move and attack but is weak.", "👮‍♀️", 4, 6)
	// Major is a rank 7 piece
	Major = *NewPieceType("Major", '7', true, "The piece that can move and attack but is weak.", "👮‍♂️", 3, 7)
	// Colonel is a rank 8 piece
	Colonel = *NewPieceType("Colonel", '8', true, "The piece that can move and attack but is weak.", "👮‍♀️", 2, 8)
	// General is a rank 9 piece
	General = *NewPieceType("General", '9', true, "The piece that can move and attack but is weak.", "👮‍♂️", 1, 9)
	// Marshal is the highest rank piece (10)
	Marshal = *NewPieceType("Marshal", 'M', true, "The piece that can move and attack but is weak.", "👮‍♀️", 1, 10)
)
