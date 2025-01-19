package entities

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Model            rl.Model
	Position         rl.Vector3
	ForwardPosition  rl.Vector3
	AimPosition      rl.Vector3
	Rotation         rl.Vector3
	Movement         rl.Vector3 // should be normalized
	WalkingSpeed     float32
	StrafingSpeed    float32
	JumpingSpeed     float32
	DashSpeed        float32
	VerticalMovement float32
	DashModifier     float32
	DashDirection    rl.Vector3
	Camera           *rl.Camera3D
	CameraSpeed      float32
}

func NewPlayer() Player {
	model := rl.LoadModelFromMesh(rl.GenMeshCube(1, 3, 1))
	position := rl.NewVector3(0, 2, 0)
	forward := rl.NewVector3(0, 0, 0)
	aimpos := rl.NewVector3(0, 0, 0)
	rotation := rl.NewVector3(0, 0, 0)

	camera := rl.NewCamera3D(
		rl.NewVector3(-20.0, 10.0, 0.0),
		position,
		rl.NewVector3(0, 1.0, 0),
		45.0,
		rl.CameraPerspective,
	)

	return Player{
		Model:           model,
		Position:        position,
		ForwardPosition: forward,
		AimPosition:     aimpos,
		Rotation:        rotation,
		WalkingSpeed:    20.0,
		StrafingSpeed:   10,
		JumpingSpeed:    80,
		DashSpeed:       160,
		Camera:          &camera,
		CameraSpeed:     100,
	}
}

func (Player *Player) ForwardDirection() rl.Vector3 {
	return rl.Vector3Normalize(rl.Vector3Subtract(Player.ForwardPosition, Player.Position))
}

func (Player *Player) RightDirection() rl.Vector3 {
	return rl.Vector3CrossProduct(Player.ForwardDirection(), rl.NewVector3(0, 1, 0))
}

func (Player *Player) MoveByVector(vector rl.Vector3, speedmod float32) {
	speed := Player.WalkingSpeed * rl.GetFrameTime() * speedmod
	vector = rl.Vector3Multiply(vector, rl.NewVector3(speed, speed, speed))

	Player.Position = rl.Vector3Add(Player.Position, vector)
	Player.AimPosition = rl.Vector3Add(Player.AimPosition, vector)
	Player.ForwardPosition = rl.Vector3Add(Player.ForwardPosition, vector)
	Player.Camera.Position = rl.Vector3Add(Player.Camera.Position, vector)
	Player.Camera.Target = rl.Vector3Add(Player.Camera.Target, vector)

	Player.AimPosition.Y = Player.Position.Y + 2
	Player.ForwardPosition.Y = Player.Position.Y
}

func (Player *Player) Jump() {
	if Player.VerticalMovement > 0 || Player.Position.Y > 2 {
		return
	}
	Player.VerticalMovement = Player.JumpingSpeed
}

func (Player *Player) Dash(direction rl.Vector3) {
	if Player.DashModifier > 0 || Player.Position.Y > 2 {
		return
	}
	Player.DashModifier = Player.DashSpeed / Player.WalkingSpeed
	Player.DashDirection = direction
}

func (Player *Player) GravityAndPositionLoop() {
	if Player.VerticalMovement > 0 {
		Player.MoveByVector(rl.NewVector3(0, Player.VerticalMovement*rl.GetFrameTime(), 0), 1)
		Player.VerticalMovement -= 100 * rl.GetFrameTime()
	}
	if Player.Position.Y > 2 {
		Player.MoveByVector(rl.NewVector3(0, Player.VerticalMovement*rl.GetFrameTime(), 0), 1)
		Player.VerticalMovement -= 100 * rl.GetFrameTime()
	}
	if Player.DashModifier > 0 {
		Player.MoveByVector(Player.DashDirection, Player.DashModifier)
		Player.DashModifier -= Player.DashSpeed / Player.WalkingSpeed * rl.GetFrameTime() * 3
	} else {
		Player.MoveByVector(Player.Movement, 1)
	}

}

func (Player *Player) RenderHud() {

}
