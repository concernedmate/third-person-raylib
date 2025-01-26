package main

import (
	"concernedmate/trial-raylib/controls"
	"concernedmate/trial-raylib/entities"
	"concernedmate/trial-raylib/gameplay"
	"concernedmate/trial-raylib/hud"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1280, 720, "monhun bow clone")
	defer rl.CloseWindow()

	World := gameplay.World{
		MainPlayer: entities.NewPlayer(),
	}

	World.Mobs = append(World.Mobs, entities.NewUndead(rl.NewVector3(10, 1, 0)))

	rl.DisableCursor()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.ClearBackground(rl.White)

		// controls
		controls.UpdateCameraThirdPerson(&World.MainPlayer)
		controls.UpdatePlayerMovement(&World.MainPlayer)
		controls.UpdateChargeLevel(&World.MainPlayer, &World)

		// game state
		World.LoopPhysicsEntities()
		World.LoopGarbageDeletionEntities()

		// start draw
		rl.BeginDrawing()

		// render
		World.LoopRenderEntities()

		// HUD
		hud.RenderHud(&World.MainPlayer)

		rl.EndDrawing()
	}
}
