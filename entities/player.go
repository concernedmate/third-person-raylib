package entities

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const PLAYER_DASH_CD = 600 // ms

type Player struct {
	Model           rl.Model
	Position        rl.Vector3
	ForwardPosition rl.Vector3
	AimPosition     rl.Vector3
	Rotation        rl.Vector3
	Movement        rl.Vector3

	WalkingSpeed    float32
	WalkingModifier float32

	JumpingSpeed     float32
	VerticalMovement float32
	FallingSpeed     float32

	DashSpeed     float32
	DashModifier  float32
	DashDirection rl.Vector3
	DashTimer     time.Time

	Bow                Bow
	ChargeSpeed        int
	ChargeCurrentLevel float32
	ChargeLevel1       int
	ChargeLevel2       int
	ChargeLevel3       int
	ChargeTimer        time.Time

	Camera      *rl.Camera3D
	CameraSpeed float32
}

func NewPlayer() Player {
	model := rl.LoadModelFromMesh(rl.GenMeshCube(0.75, 2, 0.75))
	position := rl.NewVector3(0, 1, 0)
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

	bow := NewBowType1()

	return Player{
		Model:           model,
		Position:        position,
		ForwardPosition: forward,
		AimPosition:     aimpos,
		Rotation:        rotation,

		WalkingSpeed:    10.0,
		WalkingModifier: 0.2,
		JumpingSpeed:    30,
		DashSpeed:       50,

		Bow:                bow,
		ChargeSpeed:        75,
		ChargeCurrentLevel: 0,
		ChargeLevel1:       100,
		ChargeLevel2:       200,
		ChargeLevel3:       300,

		Camera:      &camera,
		CameraSpeed: 100,
	}
}

func (Player *Player) ForwardDirection() rl.Vector3 {
	return rl.Vector3Normalize(rl.Vector3Subtract(Player.ForwardPosition, Player.Position))
}

func (Player *Player) RightDirection() rl.Vector3 {
	return rl.Vector3CrossProduct(Player.ForwardDirection(), rl.NewVector3(0, 1, 0))
}

func (Player *Player) AimForwardDirection() rl.Vector3 {
	return rl.Vector3Normalize(rl.Vector3Subtract(Player.Camera.Target, Player.AimPosition))
}

// vector must be normalized,
// speedmod is modifier that is applied to WalkingSpeed
func (Player *Player) MoveByVector(vector rl.Vector3, speedmod float32) {
	speed := Player.WalkingSpeed * rl.GetFrameTime() * speedmod
	vector = rl.Vector3Multiply(vector, rl.NewVector3(speed, speed, speed))

	Player.Position = rl.Vector3Add(Player.Position, vector)
	Player.AimPosition = rl.Vector3Add(Player.AimPosition, vector)
	Player.ForwardPosition = rl.Vector3Add(Player.ForwardPosition, vector)
	Player.Camera.Position = rl.Vector3Add(Player.Camera.Position, vector)
	Player.Camera.Target = rl.Vector3Add(Player.Camera.Target, vector)

	Player.AimPosition.Y = Player.Position.Y + 1.3
	Player.ForwardPosition.Y = Player.Position.Y
}

func (Player *Player) Jump() {
	if Player.VerticalMovement > 0 || Player.Position.Y > 2 {
		return
	}
	Player.VerticalMovement = Player.JumpingSpeed
}

func (Player *Player) Dash(direction rl.Vector3) {
	if time.Since(Player.DashTimer).Milliseconds() < PLAYER_DASH_CD {
		return
	}
	if Player.DashModifier > 0 || Player.Position.Y > 2 {
		return
	}
	Player.DashModifier = Player.DashSpeed / Player.WalkingSpeed
	Player.DashDirection = direction
	Player.DashTimer = time.Now()
}

func (Player *Player) Move(direction rl.Vector3, modifier float32) {
	if Player.Position.Y > 1 {
		return
	}
	Player.Movement = direction
	Player.WalkingModifier = modifier
}

func (Player *Player) ChargeArrow() {
	if Player.DashModifier > 0 {
		return
	}
	if time.Since(Player.ChargeTimer).Milliseconds() < PLAYER_DASH_CD {
		return
	}
	Player.ChargeCurrentLevel += float32(Player.ChargeSpeed) * rl.GetFrameTime()
}

// returns charge level (0, 1, 2, or 3) and reset it
func (Player *Player) ReleaseArrow() int {
	var chargeLevel int
	if Player.ChargeCurrentLevel > float32(Player.ChargeLevel1) {
		chargeLevel = 1
	}
	if Player.ChargeCurrentLevel > float32(Player.ChargeLevel2) {
		chargeLevel = 2
	}
	if Player.ChargeCurrentLevel > float32(Player.ChargeLevel3) {
		chargeLevel = 3
	}
	Player.ChargeCurrentLevel = 0
	Player.ChargeTimer = time.Now()
	Player.DashTimer = time.Now()
	return chargeLevel
}

/* ================================== LOOP BELOW ================================== */
/* ================================== LOOP BELOW ================================== */

func (Player *Player) GravityAndPositionLoop() {
	if Player.VerticalMovement > 0 {
		Player.MoveByVector(rl.NewVector3(0, Player.VerticalMovement*rl.GetFrameTime(), 0), 1)
	}
	if Player.Position.Y > 1 {
		Player.MoveByVector(rl.NewVector3(0, Player.VerticalMovement*rl.GetFrameTime(), 0), 1)
		Player.FallingSpeed = Player.FallingSpeed + 10*rl.GetFrameTime()
		Player.VerticalMovement -= Player.FallingSpeed
	} else {
		Player.FallingSpeed = 0
	}
	if Player.DashModifier > 0 {
		Player.MoveByVector(Player.DashDirection, Player.DashModifier)
		Player.DashModifier -= Player.DashSpeed / Player.WalkingSpeed * rl.GetFrameTime() * 1.65

		Player.ChargeCurrentLevel += float32(Player.ChargeSpeed) * 1.85 * rl.GetFrameTime()
	} else {
		if time.Since(Player.DashTimer).Milliseconds() < PLAYER_DASH_CD {
			return
		}
		Player.MoveByVector(Player.Movement, Player.WalkingModifier)
	}
}
