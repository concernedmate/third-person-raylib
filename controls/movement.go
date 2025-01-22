package controls

import (
	"concernedmate/trial-raylib/entities"
	"concernedmate/trial-raylib/gameplay"
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
	if rl.IsKeyPressed(rl.KeySpace) {
		player.Dash(vectorMovement)
	}
	if rl.IsKeyPressed(rl.KeyC) {
		player.Jump()
	}

	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		player.Move(vectorMovement, 0.3)
	} else {
		player.Move(vectorMovement, 1)
	}

	diff := rl.Vector3Subtract(player.Position, player.Camera.Position)
	yAngle := math.Atan2(float64(diff.Z), float64(diff.X)) + math.Pi/2.0
	rotation := rl.NewVector3(0, float32(yAngle), 0)

	player.Rotation = rotation
	player.Model.Transform = rl.MatrixRotateY(rotation.Y)

	player.ForwardPosition.X = player.Position.X + float32(math.Sin(yAngle))*10
	player.ForwardPosition.Z = player.Position.Z - float32(math.Cos(yAngle))*10

	aimPos := rl.Vector3Subtract(player.Position, rl.Vector3Multiply(player.ForwardDirection(), rl.NewVector3(3, 3, 3)))
	player.AimPosition.X = aimPos.X + float32(math.Sin(yAngle+(90*math.Pi/180)))
	player.AimPosition.Z = aimPos.Z - float32(math.Cos(yAngle+(90*math.Pi/180)))
}

func UpdateCameraThirdPerson(player *entities.Player) {
	mouseDelta := rl.GetMouseDelta()

	camSpeed := player.CameraSpeed * rl.GetFrameTime()

	var dist float32
	var diff rl.Vector3
	if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		player.Camera.Position = player.AimPosition
	}
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		AimingCameraYaw(player, -mouseDelta.X*0.001)
		AimingCameraPitch(player, -mouseDelta.Y*0.001)
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

func AimingCameraYaw(player *entities.Player, angle float32) {
	camera := player.Camera
	up := rl.GetCameraUp(camera)

	cameraPosition := rl.Vector3Subtract(player.Position, camera.Position)

	cameraPosition = rl.Vector3RotateByAxisAngle(cameraPosition, up, angle)
	camera.Position = rl.Vector3Subtract(player.Position, cameraPosition)
}

func AimingCameraPitch(player *entities.Player, angle float32) {
	camera := player.Camera

	up := rl.GetCameraUp(camera)
	right := rl.GetCameraRight(camera)

	// View vector
	tempPosition := player.Position
	tempPosition.Y += 1.25
	cameraPosition := rl.Vector3Subtract(tempPosition, camera.Position)

	// Clamp view up
	maxAngleUp := rl.Vector3Angle(up, cameraPosition)
	maxAngleUp = maxAngleUp - 0.001 // avoid numerical errors
	if angle > maxAngleUp-1 {
		angle = maxAngleUp - 1
	}
	// Clamp view down
	maxAngleDown := rl.Vector3Angle(rl.Vector3Negate(up), cameraPosition)
	maxAngleDown = maxAngleDown * -1.0  // downwards angle is negative
	maxAngleDown = maxAngleDown + 0.001 // avoid numerical errors
	if angle < maxAngleDown+1.3 {
		angle = maxAngleDown + 1.3
	}

	// Rotate view vector around right axis
	cameraPosition = rl.Vector3RotateByAxisAngle(cameraPosition, right, angle)
	camera.Position = rl.Vector3Subtract(tempPosition, cameraPosition)

	targetPosition := rl.Vector3Subtract(tempPosition, camera.Position)
	camera.Target = rl.Vector3Add(camera.Position, rl.Vector3Multiply(targetPosition, rl.NewVector3(15/maxAngleUp, 15/maxAngleUp, 15/maxAngleUp)))
}

func ShootArrow(player *entities.Player, world *gameplay.World) {
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		player.ChargeArrow()
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			chargeLevel := player.ReleaseArrow()
			switch chargeLevel {
			case 0:
				{
					position := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.1, 0.1, 0.1)))
					arrow := entities.NewBowProjectile(position, player.Camera.Target)
					world.BowProjectiles = append(world.BowProjectiles, arrow)
					break
				}
			case 1:
				{
					position1 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.1, 0.1, 0.1)))
					position2 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.5, 0.5, 0.5)))

					arrow1 := entities.NewBowProjectile(position1, player.Camera.Target)
					arrow2 := entities.NewBowProjectile(position2, player.Camera.Target)

					world.BowProjectiles = append(world.BowProjectiles, arrow1)
					world.BowProjectiles = append(world.BowProjectiles, arrow2)
					break
				}
			case 2:
				{
					position1 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.1, 0.1, 0.1)))
					position2 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.5, 0.5, 0.5)))
					position2.Y += 0.5
					position3 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(-0.5, -0.5, -0.5)))
					position3.Y -= 0.5

					arrow1 := entities.NewBowProjectile(position1, player.Camera.Target)
					arrow2 := entities.NewBowProjectile(position2, player.Camera.Target)
					arrow3 := entities.NewBowProjectile(position3, player.Camera.Target)

					world.BowProjectiles = append(world.BowProjectiles, arrow1)
					world.BowProjectiles = append(world.BowProjectiles, arrow2)
					world.BowProjectiles = append(world.BowProjectiles, arrow3)
					break
				}
			case 3:
				{
					position1 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.0, 0.0, 0.0)))
					position2 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.15, 0.15, 0.15)))
					position3 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(-0.15, -0.15, -0.15)))
					position4 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(0.25, 0.25, 0.25)))
					position5 := rl.Vector3Add(player.Position, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(-0.25, -0.25, -0.25)))

					target2 := rl.Vector3Add(player.Camera.Target, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(2, 2, 2)))
					target3 := rl.Vector3Add(player.Camera.Target, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(-2, -2, -2)))
					target4 := rl.Vector3Add(player.Camera.Target, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(4, 4, 4)))
					target5 := rl.Vector3Add(player.Camera.Target, rl.Vector3Multiply(player.RightDirection(), rl.NewVector3(-4, -4, -4)))

					arrow1 := entities.NewBowProjectile(position1, player.Camera.Target)
					arrow2 := entities.NewBowProjectile(position2, target2)
					arrow3 := entities.NewBowProjectile(position3, target3)
					arrow4 := entities.NewBowProjectile(position4, target4)
					arrow5 := entities.NewBowProjectile(position5, target5)

					world.BowProjectiles = append(world.BowProjectiles, arrow1)
					world.BowProjectiles = append(world.BowProjectiles, arrow2)
					world.BowProjectiles = append(world.BowProjectiles, arrow3)
					world.BowProjectiles = append(world.BowProjectiles, arrow4)
					world.BowProjectiles = append(world.BowProjectiles, arrow5)
					break
				}
			}
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseButtonRight) {
		player.ReleaseArrow()
	}
}
