package entities

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BowProjectile struct {
	Model    rl.Model
	Position rl.Vector3
	Target   rl.Vector3

	Speed            float32
	VerticalMovement float32

	LifeTimer      time.Time
	LifeDurationMS int

	MarkForDeletion bool
}

func NewBowProjectile(position rl.Vector3, target rl.Vector3) BowProjectile {
	// rl.GenMeshCone(0.1, 1.5, 10)
	model := rl.LoadModelFromMesh(rl.GenMeshSphere(0.2, 10, 10))

	diff := rl.Vector3Subtract(target, position)
	targetNormalized := rl.Vector3Normalize(diff)

	verticalMovement := 3

	entity := BowProjectile{
		Model:            model,
		Position:         position,
		Target:           targetNormalized,
		VerticalMovement: float32(verticalMovement),
		Speed:            100,
		LifeTimer:        time.Now(),
		LifeDurationMS:   2000,
		MarkForDeletion:  false,
	}
	return entity
}

// vector normalized
func (Arrow *BowProjectile) MoveByVector(vector rl.Vector3, speedmod float32) {
	speed := Arrow.Speed * rl.GetFrameTime() * speedmod
	vector = rl.Vector3Multiply(vector, rl.NewVector3(speed, speed, speed))
	Arrow.Position = rl.Vector3Add(Arrow.Position, vector)
}

func (Arrow *BowProjectile) GravityAndPositionLoop() {
	if time.Since(Arrow.LifeTimer).Milliseconds() > int64(Arrow.LifeDurationMS) {
		Arrow.MarkForDeletion = true
		return
	}
	if Arrow.Position.Y > 0 {
		Arrow.MoveByVector(rl.NewVector3(0, Arrow.VerticalMovement*rl.GetFrameTime(), 0), 1)
		Arrow.VerticalMovement -= 10 * rl.GetFrameTime()

		Arrow.MoveByVector(Arrow.Target, 1)
	}

}
