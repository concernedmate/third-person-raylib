package entities

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Mob struct {
	Model       rl.Model
	Position    rl.Vector3
	AggroTarget *Player

	MovementSpeed int
	Health        int

	Damage int
	Armor  int

	MarkForDeletion bool
}

type Projectile struct {
	Model    rl.Model
	Position rl.Vector3
	Target   rl.Vector3

	Speed            float32
	VerticalMovement float32

	LifeTimer      time.Time
	LifeDurationMS int

	MarkForDeletion bool
}

// vector normalized
func (Proj *Projectile) MoveByVector(vector rl.Vector3, speedmod float32) {
	speed := Proj.Speed * rl.GetFrameTime() * speedmod
	vector = rl.Vector3Multiply(vector, rl.NewVector3(speed, speed, speed))
	Proj.Position = rl.Vector3Add(Proj.Position, vector)
}

func (Proj *Projectile) GravityAndPositionLoop() {
	if time.Since(Proj.LifeTimer).Milliseconds() > int64(Proj.LifeDurationMS) {
		Proj.MarkForDeletion = true
		return
	}
	if Proj.Position.Y > 0 {
		Proj.MoveByVector(rl.NewVector3(0, Proj.VerticalMovement*rl.GetFrameTime(), 0), 1)
		Proj.VerticalMovement -= 10 * rl.GetFrameTime()

		Proj.MoveByVector(Proj.Target, 1)
	}
}

func (Proj *Projectile) FreezePosition() {
	Proj.Speed = 0
	Proj.VerticalMovement = 0
}
