package main

import (
	"concernedmate/trial-raylib/controls"
	"concernedmate/trial-raylib/entities"
	"concernedmate/trial-raylib/gameplay"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1280, 720, "monhun bow clone")
	defer rl.CloseWindow()

	World := gameplay.World{
		MainPlayer: entities.NewPlayer(),
	}

	rl.DisableCursor()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.ClearBackground(rl.White)

		// controls
		controls.UpdateCameraThirdPerson(&World.MainPlayer)
		controls.UpdatePlayerMovement(&World.MainPlayer)
		controls.ShootArrow(&World.MainPlayer, &World)

		// game state
		World.LoopPhysicsEntities()
		World.LoopGarbageDeletionEntities()

		// start draw
		rl.BeginDrawing()

		// render
		World.RenderEntities()

		// HUD
		World.MainPlayer.RenderHud()

		rl.EndDrawing()
	}
}
