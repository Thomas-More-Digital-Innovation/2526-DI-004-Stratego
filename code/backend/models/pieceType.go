package models

type PieceType struct {
	name           string
	rank           byte
	movable        bool
	description    string
	icon           string
	count          int
	strategicValue int
}

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

func (pt *PieceType) GetName() string {
	return pt.name
}

func (pt *PieceType) GetRank() byte {
	return pt.rank
}

func (pt *PieceType) IsMovable() bool {
	return pt.movable
}
func (pt *PieceType) GetDescription() string {
	return pt.description
}

func (pt *PieceType) GetIcon() string {
	return pt.icon
}
func (pt *PieceType) GetCount() int {
	return pt.count
}

func (pt *PieceType) GetStrategicValue() int {
	return pt.strategicValue
}
