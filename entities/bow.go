package entities

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BowType int

const (
	Focus BowType = iota
	Spread
	Pierce
	Blast
)

type Bow struct {
	L1Type      BowType
	L2Type      BowType
	L3Type      BowType
	L1ProjCount int
	L2ProjCount int
	L3ProjCount int
}

func NewBowType1() Bow {
	return Bow{
		L1Type:      Focus,
		L2Type:      Focus,
		L3Type:      Spread,
		L1ProjCount: 2,
		L2ProjCount: 3,
		L3ProjCount: 6,
	}
}

func NewBowProjectile(position rl.Vector3, target rl.Vector3) Projectile {
	// rl.GenMeshCone(0.1, 1.5, 10)
	model := rl.LoadModelFromMesh(rl.GenMeshSphere(0.1, 10, 10))

	diff := rl.Vector3Subtract(target, position)
	targetNormalized := rl.Vector3Normalize(diff)

	verticalMovement := 3

	entity := Projectile{
		Model:            model,
		Position:         position,
		Target:           targetNormalized,
		VerticalMovement: float32(verticalMovement),
		Speed:            50,
		LifeTimer:        time.Now(),
		LifeDurationMS:   2000,
		MarkForDeletion:  false,
	}
	return entity
}
