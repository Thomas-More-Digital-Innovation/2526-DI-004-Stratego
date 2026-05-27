package models

// PieceType defines the characteristics of a game piece
type PieceType struct {
	name           string
	rank           byte
	movable        bool
	description    string
	icon           string
	count          int
	strategicValue int
}

// NewPieceType creates a new PieceType instance
func NewPieceType(name string, rank byte, movable bool, description string, icon string, count int, strategicValue int) *PieceType {
	return &PieceType{
		name:           name,
		rank:           rank,
		movable:        movable,
		description:    description,
		icon:           icon,
		count:          count,
		strategicValue: strategicValue,
	}
}

// GetName returns the name of the piece type
func (pt *PieceType) GetName() string {
	return pt.name
}

// GetRank returns the rank of the piece type
func (pt *PieceType) GetRank() byte {
	return pt.rank
}

// IsMovable returns true if the piece type is movable
func (pt *PieceType) IsMovable() bool {
	return pt.movable
}

// GetDescription returns the description of the piece type
func (pt *PieceType) GetDescription() string {
	return pt.description
}

// GetIcon returns the icon of the piece type
func (pt *PieceType) GetIcon() string {
	return pt.icon
}

// GetCount returns the initial count of this piece type in a game
func (pt *PieceType) GetCount() int {
	return pt.count
}

// GetStrategicValue returns the strategic value of the piece type for AI
func (pt *PieceType) GetStrategicValue() int {
	return pt.strategicValue
}
