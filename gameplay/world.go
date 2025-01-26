package gameplay

import (
	"concernedmate/trial-raylib/entities"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var model *rl.Model

type World struct {
	MainPlayer   entities.Player
	OtherPlayers []entities.Player // only for rendering

	Projectiles []entities.Projectile
	Mobs        []entities.Mob
}

func (World *World) renderMobs(mob entities.Mob) {
	rl.DrawModel(mob.Model, mob.Position, 1, rl.Red)
}

func (World *World) checkMobsProjCollision(mob entities.Mob, proj entities.Projectile) bool {
	a := rl.GetModelBoundingBox(mob.Model)
	a.Min = rl.Vector3Add(a.Min, mob.Position)
	a.Max = rl.Vector3Add(a.Max, mob.Position)

	b := rl.GetModelBoundingBox(proj.Model)
	b.Min = rl.Vector3Add(b.Min, proj.Position)
	b.Max = rl.Vector3Add(b.Max, proj.Position)

	return rl.CheckCollisionBoxes(a, b)
}

/* ================================== LOOP BELOW ================================== */
/* ================================== LOOP BELOW ================================== */

func (World *World) LoopPhysicsEntities() {
	World.MainPlayer.GravityAndPositionLoop()

	for idx := range World.Projectiles {
		World.Projectiles[idx].GravityAndPositionLoop()
	}

	// collision
	for _, mob := range World.Mobs {
		for idx, proj := range World.Projectiles {
			if World.checkMobsProjCollision(mob, proj) {
				World.Projectiles[idx].FreezePosition()
			}
		}
	}
}

func (World *World) LoopGarbageDeletionEntities() {
	var newProj []entities.Projectile
	for _, val := range World.Projectiles {
		if !val.MarkForDeletion {
			newProj = append(newProj, val)
		}
	}
	World.Projectiles = newProj

	var newMobs []entities.Mob
	for _, val := range World.Mobs {
		if !val.MarkForDeletion {
			newMobs = append(newMobs, val)
		}
	}
	World.Mobs = newMobs
}

func (World *World) LoopRenderEntities() {
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

	// render
	rl.DrawModelWires(World.MainPlayer.Model, World.MainPlayer.Position, 1, rl.Blue)
	for _, proj := range World.Projectiles {
		rl.DrawModel(proj.Model, proj.Position, 1, rl.Green)
	}

	for _, mob := range World.Mobs {
		World.renderMobs(mob)
	}

	rl.DrawGrid(100, 1.0)
	rl.EndMode3D()
}
