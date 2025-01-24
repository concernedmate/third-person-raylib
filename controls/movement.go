package controls

import (
	"concernedmate/trial-raylib/entities"
	"concernedmate/trial-raylib/gameplay"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func aimingCameraYaw(Player *entities.Player, angle float32) {
	camera := Player.Camera
	up := rl.GetCameraUp(camera)

	cameraPosition := rl.Vector3Subtract(Player.Position, camera.Position)

	cameraPosition = rl.Vector3RotateByAxisAngle(cameraPosition, up, angle)
	camera.Position = rl.Vector3Subtract(Player.Position, cameraPosition)
}

func aimingCameraPitch(Player *entities.Player, angle float32) {
	camera := Player.Camera

	up := rl.GetCameraUp(camera)
	right := rl.GetCameraRight(camera)

	// View vector
	tempPosition := Player.Position
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

func releaseArrow(Player *entities.Player, World *gameplay.World, bowType entities.BowType, projCount int) {
	switch bowType {
	case entities.Focus:
		{
			for i := range projCount {
				positionModifier := rand.Float32() / 3
				if i%2 == 0 {
					positionModifier *= -1
				}
				vectorModifier := rl.Vector3Multiply(Player.RightDirection(), rl.NewVector3(positionModifier, positionModifier, positionModifier))
				if i%2 == 0 {
					vectorModifier.Y += 0.5
				}
				position := rl.Vector3Add(Player.Position, vectorModifier)
				arrow := entities.NewBowProjectile(position, Player.Camera.Target)
				World.BowProjectiles = append(World.BowProjectiles, arrow)
			}
			break
		}
	case entities.Spread:
		{
			for i := range projCount {
				positionModifier := rand.Float32() / 3
				targetPositionModifier := 1 + (positionModifier * 10)

				if i%2 == 0 {
					positionModifier *= -1
					targetPositionModifier *= -1
				}

				vectorModifier := rl.Vector3Multiply(Player.RightDirection(), rl.NewVector3(positionModifier, positionModifier, positionModifier))
				targetVectorModifier := rl.Vector3Multiply(Player.RightDirection(), rl.NewVector3(targetPositionModifier, targetPositionModifier, targetPositionModifier))

				position := rl.Vector3Add(Player.Position, vectorModifier)
				targetPosition := rl.Vector3Add(Player.Camera.Target, targetVectorModifier)

				arrow := entities.NewBowProjectile(position, targetPosition)
				World.BowProjectiles = append(World.BowProjectiles, arrow)
			}
			break
		}
	}
}

func UpdatePlayerMovement(Player *entities.Player) {
	// calculate position
	forward := Player.ForwardDirection()
	right := Player.RightDirection()

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
		Player.Dash(vectorMovement)
	}
	if rl.IsKeyPressed(rl.KeyC) {
		Player.Jump()
	}

	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		Player.Move(vectorMovement, 0.3)
	} else {
		Player.Move(vectorMovement, 1)
	}

	diff := rl.Vector3Subtract(Player.Position, Player.Camera.Position)
	yAngle := math.Atan2(float64(diff.Z), float64(diff.X)) + math.Pi/2.0
	rotation := rl.NewVector3(0, float32(yAngle), 0)

	Player.Rotation = rotation
	Player.Model.Transform = rl.MatrixRotateY(rotation.Y)

	Player.ForwardPosition.X = Player.Position.X + float32(math.Sin(yAngle))*10
	Player.ForwardPosition.Z = Player.Position.Z - float32(math.Cos(yAngle))*10

	aimPos := rl.Vector3Subtract(Player.Position, rl.Vector3Multiply(Player.ForwardDirection(), rl.NewVector3(3, 3, 3)))
	Player.AimPosition.X = aimPos.X + float32(math.Sin(yAngle+(90*math.Pi/180)))
	Player.AimPosition.Z = aimPos.Z - float32(math.Cos(yAngle+(90*math.Pi/180)))
}

func UpdateCameraThirdPerson(Player *entities.Player) {
	mouseDelta := rl.GetMouseDelta()

	camSpeed := Player.CameraSpeed * rl.GetFrameTime()

	var dist float32
	var diff rl.Vector3
	if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		Player.Camera.Position = Player.AimPosition
	}
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		aimingCameraYaw(Player, -mouseDelta.X*0.001)
		aimingCameraPitch(Player, -mouseDelta.Y*0.001)
	} else {
		rl.CameraYaw(Player.Camera, -mouseDelta.X*0.003, 1)
		rl.CameraPitch(Player.Camera, -mouseDelta.Y*0.003, 1, 1, 0)

		dist = rl.Vector3Distance(Player.Position, Player.Camera.Position)
		diff = rl.Vector3Normalize(rl.Vector3Subtract(Player.Camera.Target, Player.Camera.Position))

		if dist <= 20 {
			Player.Camera.Position = rl.Vector3Subtract(Player.Camera.Position, rl.Vector3Multiply(diff, rl.NewVector3(camSpeed, camSpeed, camSpeed)))
			Player.Camera.Target = rl.Vector3Subtract(Player.Camera.Target, rl.Vector3Multiply(diff, rl.NewVector3(camSpeed, camSpeed, camSpeed)))
		} else {
			Player.Camera.Target = Player.Position
		}
	}
}

func UpdateChargeLevel(Player *entities.Player, World *gameplay.World) {
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		Player.ChargeArrow()
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			chargeLevel := Player.ReleaseArrow()
			switch chargeLevel {
			case 1:
				{
					releaseArrow(Player, World, Player.Bow.L1Type, Player.Bow.L1ProjCount)
					break
				}
			case 2:
				{
					releaseArrow(Player, World, Player.Bow.L2Type, Player.Bow.L2ProjCount)
					break
				}
			case 3:
				{
					releaseArrow(Player, World, Player.Bow.L3Type, Player.Bow.L3ProjCount)
					break
				}
			}
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseButtonRight) {
		Player.ReleaseArrow()
	}
}
