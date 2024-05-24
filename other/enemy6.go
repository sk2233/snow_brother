package other

import (
	"gameBase/factory"
	"gameBase/loader"
	"gameBase/model"
	"gameBase/other"
	"gameBase/utils"
)

func init() {
	factory.RegisterPosFactory("enemy6", createEnemy6)
}

func createEnemy6(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Enemy6)
	sprites := [][]other.Sprite{
		{loader.LoadStaticSprite("item", "6/idle")},
		{loader.LoadDynamicSprite("item", "6/move")},
		{loader.LoadStaticSprite("item", "6/jump")},
		{loader.LoadDynamicSprite("item", "6/shake")},
		{loader.LoadDynamicSprite("item", "6/attack")},
	}
	res.WalkEnemy = NewWalkEnemy(res.triggerAttack, sprites)
	factory.FillPosObject(res.PosObject, o)
	res.typeName = "enemy6"
	res.speed = 1
	res.actions = append(res.actions, res.attack)
	res.dieSpriteNames = []string{"6/die/move", "6/die/idle"}
	return res
}

// Enemy6 大红球
type Enemy6 struct {
	*WalkEnemy // 4 攻击状态
}

func (e *Enemy6) attack() {
	if e.MoveHorizontal("collision", WallTag, e.ScaleX*2) {
		e.randState()
	}
}

func (e *Enemy6) triggerAttack() {
	if player != nil {
		e.ScaleX = utils.Sign(player.X - e.X)
	}
	e.switchState(4, 0)
}
