package hud

import (
	"concernedmate/trial-raylib/entities"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RenderHud(Player *entities.Player) {
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		rl.DrawCircle(640, 360, 15, rl.Fade(rl.Red, 0.5))
		rl.DrawCircle(640, 360, 2, rl.Red)
		rl.DrawCircleLines(640, 360, 15, rl.Red)
		rl.DrawCircleLines(640, 360, 10, rl.Red)

		if time.Since(Player.ChargeTimer).Milliseconds() < entities.PLAYER_DASH_CD {
			return
		}
		if Player.ChargeCurrentLevel < float32(Player.ChargeLevel1) {
			circle := 65 - (Player.ChargeCurrentLevel / float32(Player.ChargeLevel1) * 50)
			rl.DrawCircleLines(640, 360, circle, rl.Red)
		}
		if Player.ChargeCurrentLevel > float32(Player.ChargeLevel1) && Player.ChargeCurrentLevel < float32(Player.ChargeLevel2) {
			circle := 65 - ((Player.ChargeCurrentLevel - 100) / float32(Player.ChargeLevel1) * 50)
			rl.DrawCircleLines(640, 360, circle, rl.Green)
		}
		if Player.ChargeCurrentLevel > float32(Player.ChargeLevel2) && Player.ChargeCurrentLevel < float32(Player.ChargeLevel3) {
			circle := 65 - ((Player.ChargeCurrentLevel - 200) / float32(Player.ChargeLevel1) * 50)
			rl.DrawCircleLines(640, 360, circle, rl.Blue)
		}
	}
}
