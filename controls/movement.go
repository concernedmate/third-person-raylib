package controls

import (
	"concernedmate/trial-raylib/entities"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func UpdatePlayerMovement(player *entities.Player) {
	// calculate position
	forward := player.ForwardDirection()
	right := player.RightDirection()

	vectorMovement := rl.NewVector3(0, 0, 0)
	if rl.IsKeyDown(rl.KeyW) {
		vectorMovement.X += forward.X
		vectorMovement.Z += forward.Z
	}
	if rl.IsKeyDown(rl.KeyS) {
		vectorMovement.X -= forward.X
		vectorMovement.Z -= forward.Z
	}
	if rl.IsKeyDown(rl.KeyD) {
		vectorMovement.X += right.X
		vectorMovement.Z += right.Z
	}
	if rl.IsKeyDown(rl.KeyA) {
		vectorMovement.X -= right.X
		vectorMovement.Z -= right.Z
	}
	if rl.IsKeyDown(rl.KeySpace) {
		player.Dash(vectorMovement)
	}
	if rl.IsKeyDown(rl.KeyC) {
		player.Jump()
	}
	player.Movement = vectorMovement

	// calculate rotation and camera position
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		diff := rl.Vector3Subtract(player.Position, player.Camera.Target)
		yAngle := math.Atan2(float64(diff.Z), float64(diff.X)) + math.Pi/2.0
		rotation := rl.NewVector3(0, float32(yAngle), 0)

		player.Rotation = rotation
		player.Model.Transform = rl.MatrixRotateY(rotation.Y)

		player.Position.X = player.AimPosition.X + float32(math.Sin(yAngle+(130*math.Pi/180)))
		player.Position.Z = player.AimPosition.Z - float32(math.Cos(yAngle+(130*math.Pi/180)))

		player.ForwardPosition.X = player.Position.X - float32(math.Sin(yAngle))*10
		player.ForwardPosition.Z = player.Position.Z + float32(math.Cos(yAngle))*10
	} else {
		diff := rl.Vector3Subtract(player.Position, player.Camera.Position)
		yAngle := math.Atan2(float64(diff.Z), float64(diff.X)) + math.Pi/2.0
		rotation := rl.NewVector3(0, float32(yAngle), 0)

		player.Rotation = rotation
		player.Model.Transform = rl.MatrixRotateY(rotation.Y)

		player.ForwardPosition.X = player.Position.X + float32(math.Sin(yAngle))*10
		player.ForwardPosition.Z = player.Position.Z - float32(math.Cos(yAngle))*10

		player.AimPosition.X = player.Position.X + float32(math.Sin(yAngle+(130*math.Pi/180)))
		player.AimPosition.Z = player.Position.Z - float32(math.Cos(yAngle+(130*math.Pi/180)))
	}
}

func UpdateCameraThirdPerson(player *entities.Player) {
	mouseDelta := rl.GetMouseDelta()

	camSpeed := player.CameraSpeed * rl.GetFrameTime()

	var dist float32
	var diff rl.Vector3
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		rl.CameraYaw(player.Camera, -mouseDelta.X*0.003, 0)
		rl.CameraPitch(player.Camera, -mouseDelta.Y*0.003, 1, 0, 0)

		dist = rl.Vector3Distance(player.AimPosition, player.Camera.Position)
		diff = rl.Vector3Normalize(rl.Vector3Subtract(player.AimPosition, player.Camera.Position))

		if dist >= 1 {
			newPos := rl.Vector3Add(player.Camera.Position, rl.Vector3Multiply(diff, rl.NewVector3(camSpeed, camSpeed, camSpeed)))
			newPos2 := rl.Vector3Add(player.Camera.Target, rl.Vector3Multiply(diff, rl.NewVector3(camSpeed, camSpeed, camSpeed)))

			if rl.Vector3Distance(newPos, player.AimPosition) > 3 {
				newPos.Y = player.Position.Y
				player.Camera.Position = newPos
				player.Camera.Target = newPos2
			} else {
				player.Camera.Position = player.AimPosition
			}
		}
	} else {
		rl.CameraYaw(player.Camera, -mouseDelta.X*0.003, 1)
		rl.CameraPitch(player.Camera, -mouseDelta.Y*0.003, 1, 1, 0)

		dist = rl.Vector3Distance(player.Position, player.Camera.Position)
		diff = rl.Vector3Normalize(rl.Vector3Subtract(player.Camera.Target, player.Camera.Position))

		if dist <= 20 {
			player.Camera.Position = rl.Vector3Subtract(player.Camera.Position, rl.Vector3Multiply(diff, rl.NewVector3(camSpeed, camSpeed, camSpeed)))
			player.Camera.Target = rl.Vector3Subtract(player.Camera.Target, rl.Vector3Multiply(diff, rl.NewVector3(camSpeed, camSpeed, camSpeed)))
		} else {
			player.Camera.Target = player.Position
		}
	}
}
