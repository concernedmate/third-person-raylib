package gameplay

import (
	"concernedmate/trial-raylib/entities"
	"math"
	"reflect"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var model *rl.Model

type World struct {
	MainPlayer   entities.Player
	OtherPlayers []entities.Player // only for rendering

	BowProjectiles []entities.BowProjectile
}

func destroyEntities[V entities.BowProjectile](World *World, entity V) {
	v := reflect.ValueOf(entity)

	switch v.Type().Name() {
	case "BowProjectile":
		{
			var newArr []entities.BowProjectile
			for _, val := range World.BowProjectiles {
				if !val.MarkForDeletion {
					newArr = append(newArr, val)
				}
			}
			World.BowProjectiles = newArr
		}
	}
}

func (World *World) LoopPhysicsEntities() {
	World.MainPlayer.GravityAndPositionLoop()
	for idx := range World.BowProjectiles {
		World.BowProjectiles[idx].GravityAndPositionLoop()
	}
}

func (World *World) LoopGarbageDeletionEntities() {
	for _, val := range World.BowProjectiles {
		destroyEntities(World, val)
	}
}

func (World *World) RenderEntities() {
	if model == nil {
		// rl.GenMeshCone(0.1, 1.5, 10)
		m := rl.LoadModelFromMesh(rl.GenMeshCube(0.7, 0.7, 0.7))
		model = &m
	} else {
		diff := rl.Vector3Subtract(World.MainPlayer.Camera.Target, World.MainPlayer.Position)

		// TODO wtf is this
		xAngle := math.Atan2(float64(diff.Y), float64(diff.Z)) + math.Pi/2.0
		yAngle := math.Atan2(float64(diff.Z), float64(diff.X)) + math.Pi/2.0
		zAngle := math.Atan2(float64(diff.Z), float64(diff.X)) + math.Pi/2.0

		rotation := rl.NewVector3(float32(xAngle), float32(yAngle), float32(zAngle))
		model.Transform = rl.MatrixRotateXYZ(rotation)

	}
	rl.BeginMode3D(*World.MainPlayer.Camera)

	// rl.DrawModel(*model, World.MainPlayer.Camera.Target, 1, rl.LightGray)
	rl.DrawModelWires(*model, World.MainPlayer.ForwardPosition, 1, rl.Green)
	rl.DrawModelWires(*model, World.MainPlayer.Camera.Target, 1, rl.Red)

	rl.DrawModelWires(World.MainPlayer.Model, World.MainPlayer.Position, 1, rl.Blue)
	for _, arrow := range World.BowProjectiles {
		rl.DrawModel(arrow.Model, arrow.Position, 1, rl.Green)
	}

	rl.DrawGrid(100, 1.0)
	rl.EndMode3D()
}
