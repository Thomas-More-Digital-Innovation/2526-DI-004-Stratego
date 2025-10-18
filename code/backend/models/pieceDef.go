package models

var Flag PieceType = *NewPieceType("Flag", '0', false, "The piece you must capture to win the game.", "🚩", 1, 0)
var Bomb PieceType = *NewPieceType("Bomb", 'B', false, "The piece that cannot move and eliminates most attackers.", "💣", 6, 7)
var Spy PieceType = *NewPieceType("Spy", '1', true, "The piece that can move and attack but is weak.", "🕵️", 1, 7)
var Scout PieceType = *NewPieceType("Scout", '2', true, "The piece that can move multiple spaces and attack.", "🕵️‍♂️", 8, 3)
var Miner PieceType = *NewPieceType("Miner", '3', true, "The piece that can move and attack but is weak.", "⛏️", 5, 6)
var Sergeant PieceType = *NewPieceType("Sergeant", '4', true, "The piece that can move and attack but is weak.", "👮", 4, 4)
var Lieutenant PieceType = *NewPieceType("Lieutenant", '5', true, "The piece that can move and attack but is weak.", "👮‍♂️", 4, 5)
var Captain PieceType = *NewPieceType("Captain", '6', true, "The piece that can move and attack but is weak.", "👮‍♀️", 4, 6)
var Major PieceType = *NewPieceType("Major", '7', true, "The piece that can move and attack but is weak.", "👮‍♂️", 3, 7)
var Colonel PieceType = *NewPieceType("Colonel", '8', true, "The piece that can move and attack but is weak.", "👮‍♀️", 2, 8)
var General PieceType = *NewPieceType("General", '9', true, "The piece that can move and attack but is weak.", "👮‍♂️", 1, 9)
var Marshal PieceType = *NewPieceType("Marshal", 'M', true, "The piece that can move and attack but is weak.", "👮‍♀️", 1, 10)
