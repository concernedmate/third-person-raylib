package gameplay

import "concernedmate/trial-raylib/entities"

type Game struct {
	MainPlayer  entities.Player
	OtherPlayer []entities.Player // only for rendering
}

func (Game *Game) GameLoop() {
}
