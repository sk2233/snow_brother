package other

import (
	"gameBase/factory"
	"gameBase/loader"
	"gameBase/model"
	"gameBase/other"
	"gameBase/utils"
)

func init() {
	factory.RegisterPosFactory("enemy1", createEnemy1)
}

func createEnemy1(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Enemy1)
	sprites := [][]other.Sprite{
		{loader.LoadStaticSprite("item", "enemy/1/idle")},
		{loader.LoadDynamicSprite("item", "enemy/1/move")},
		{loader.LoadStaticSprite("item", "enemy/1/jump/up"), // 统一先up(都有up的)
			loader.LoadStaticSprite("item", "enemy/1/jump/down")},
		{loader.LoadDynamicSprite("item", "enemy/1/shake")},
		{loader.LoadDynamicSprite("item", "enemy/1/roll")},
	}
	res.WalkEnemy = NewWalkEnemy(res.triggerRoll, sprites)
	factory.FillPosObject(res.PosObject, o)
	res.typeName = "enemy1"
	res.speed = 1
	res.actions = append(res.actions, res.roll)
	res.dieSpriteNames = []string{"enemy/1/die/move", "enemy/1/die/idle"}
	return res
}

// Enemy1 小红球
type Enemy1 struct {
	*WalkEnemy // 4 翻滚状态
	rollTimer  int
}

func (e *Enemy1) roll() {
	e.rollTimer--
	if e.rollTimer < 0 {
		e.randState()
	}
	e.MoveHorizontal("collision", WallTag, e.ScaleX*2)
	e.checkAir()
}

func (e *Enemy1) triggerRoll() {
	if player != nil {
		e.ScaleX = utils.Sign(player.X - e.X)
	}
	e.rollTimer = 60
	e.switchState(4, 0)
}
