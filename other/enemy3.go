package other

import (
	"gameBase/factory"
	"gameBase/loader"
	"gameBase/model"
	"gameBase/other"
	"gameBase/tool"
)

func init() {
	factory.RegisterPosFactory("enemy3", createEnemy3)
}

func createEnemy3(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Enemy3)
	sprites := [][]other.Sprite{
		{loader.LoadStaticSprite("item", "3/idle")},
		{loader.LoadDynamicSprite("item", "3/move")},
		{loader.LoadStaticSprite("item", "3/jump")},
		{loader.LoadDynamicSprite("item", "3/shake")},
		{loader.LoadStaticSprite("item", "3/down")},
	}
	res.WalkEnemy = NewWalkEnemy(res.triggerThrough, sprites)
	factory.FillPosObject(res.PosObject, o)
	res.typeName = "enemy3"
	res.speed = 1
	res.actions = append(res.actions, res.through)
	res.dieSpriteNames = []string{"3/die/move", "3/die/idle"}
	return res
}

// Enemy3 穿墙怪
type Enemy3 struct {
	*WalkEnemy
	throughTimer int // 吊着动画持续时间
}

func (e *Enemy3) triggerThrough() { // 必须当前允许穿墙才能穿墙
	if tool.CollisionPoint("collision", FloorTag, e.X, e.GetBottom()+1) != nil {
		e.Y += 24
		e.throughTimer = 30
		e.switchState(4, 0)
	} else {
		e.randState()
	}
}

func (e *Enemy3) through() {
	e.throughTimer--
	if e.throughTimer < 0 {
		e.switchState(2, 0)
	}
}
