package main

import (
	"concernedmate/trial-raylib/controls"
	"concernedmate/trial-raylib/entities"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1280, 720, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	Player := entities.NewPlayer()
	// model := rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))

	rl.DisableCursor()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.ClearBackground(rl.White)

		// controls
		controls.UpdateCameraThirdPerson(&Player)
		controls.UpdatePlayerMovement(&Player)

		// player state
		Player.GravityAndPositionLoop()

		// start draw
		rl.BeginDrawing()

		// start 3d
		rl.BeginMode3D(*Player.Camera)

		rl.DrawModel(Player.Model, Player.Position, 1, rl.Blue)
		// rl.DrawModel(model, Player.ForwardPosition, 1, rl.Red)
		// rl.DrawModel(model, Player.AimPosition, 1, rl.Green)
		// rl.DrawModel(model, Player.Camera.Target, 1, rl.LightGray)
		rl.DrawGrid(100, 1.0)

		rl.EndMode3D()
		// end 3d

		rl.EndDrawing()
		// end draw
	}
}
