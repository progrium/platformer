package main

type Intent struct {
	MoveLeft  bool
	MoveRight bool
	Jump      bool
}

func EqualIntents(i1 Intent, i2 Intent) bool {
	if i1.MoveLeft == i2.MoveLeft &&
		i1.MoveRight == i2.MoveRight &&
		i1.Jump == i2.Jump {
		return true
	}
	return false
}
