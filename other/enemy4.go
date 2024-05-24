package other

import (
	"gameBase/factory"
	"gameBase/loader"
	"gameBase/model"
	"gameBase/other"
	"gameBase/sprite"
	"gameBase/tool"
	"gameBase/utils"
	"math"
)

func init() {
	factory.RegisterPosFactory("enemy4", createEnemy4)
}

func createEnemy4(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Enemy4)
	sprites := [][]other.Sprite{
		{loader.LoadStaticSprite("item", "4/idle")},
		{loader.LoadDynamicSprite("item", "4/move")},
		{loader.LoadStaticSprite("item", "4/jump/up"),
			loader.LoadStaticSprite("item", "4/jump/down")},
		{loader.LoadDynamicSprite("item", "4/shake")},
		{loader.LoadDynamicSprite("item", "4/knife")}, //飞刀
		{loader.LoadDynamicSprite("item", "4/spin/start"),
			loader.LoadDynamicSprite("item", "4/spin/loop")}, //旋风
	}
	sprites[4][0].(*sprite.DynamicSprite).AnimEnd = res.animEnd
	res.WalkEnemy = NewWalkEnemy(res.trigger, sprites)
	factory.FillPosObject(res.PosObject, o)
	res.typeName = "enemy4"
	res.speed = 1
	res.actions = append(res.actions, res.knife, res.spin)
	res.dieSpriteNames = []string{"4/die/move", "4/die/idle"}
	return res
}

// Enemy4 穿墙怪
type Enemy4 struct {
	*WalkEnemy
	spinTimer              int // 旋转持续时间
	xSpinSpeed, ySpinSpeed float64
}

func (e *Enemy4) trigger() { // 必须当前允许穿墙才能穿墙
	if player != nil {
		if math.Abs(e.Y-player.Y) < 16 { // 同一级 飞刀
			e.ScaleX = utils.Sign(player.X - e.X)
			tool.AddObject("enemy", NewKnife(e.ScaleX, e.X, e.Y-12))
			e.switchState(4, 0)
		} else {
			e.spinTimer = 300
			e.switchState(5, 0)
		}
	} else {
		e.randState()
	}
}

func (e *Enemy4) Hurt(value int) {
	if !e.inSpin() { // 旋转免疫伤害
		e.WalkEnemy.Hurt(value)
	}
}

func (e *Enemy4) knife() {

}

func (e *Enemy4) spin() {
	if e.inSpin() {
		if e.spinTimer%30 == 0 {
			e.adjustSpeed()
		}
		e.X += e.xSpinSpeed
		e.Y += e.ySpinSpeed
	} else {
		if e.spinTimer == 240 {
			e.switchState(5, 1)
			e.adjustSpeed()
		} else if e.spinTimer == 60 {
			e.switchState(5, 0)
		}
		if e.spinTimer < 0 {
			e.switchState(0, 0)
		}
	}
	e.spinTimer--
}

func (e *Enemy4) animEnd() {
	e.randState()
}

// 旋转3秒
func (e *Enemy4) inSpin() bool {
	return e.spinTimer < 240 && e.spinTimer > 60
}

func (e *Enemy4) adjustSpeed() {
	if player != nil {
		e.xSpinSpeed, e.ySpinSpeed = utils.ScaleSpeed(player.X-e.X, player.Y-e.Y, 1)
	} else { // 主角不存在 停止旋转
		e.spinTimer = 61
	}
}
