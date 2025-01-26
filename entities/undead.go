package entities

import rl "github.com/gen2brain/raylib-go/raylib"

func NewUndead(position rl.Vector3) Mob {
	model := rl.LoadModelFromMesh(rl.GenMeshCube(0.75, 2, 0.75))

	return Mob{
		Model:         model,
		Position:      position,
		Health:        100,
		MovementSpeed: 10,

		Damage: 10,
		Armor:  0,
	}
}
