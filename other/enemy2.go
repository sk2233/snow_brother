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
	factory.RegisterPosFactory("enemy2", createEnemy2)
}

func createEnemy2(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Enemy2)
	sprites := [][]other.Sprite{
		{loader.LoadStaticSprite("item", "2/idle")},
		{loader.LoadDynamicSprite("item", "2/move")},
		{loader.LoadStaticSprite("item", "2/jump")},
		{loader.LoadDynamicSprite("item", "2/shake")},
		{loader.LoadStaticSprite("item", "2/attack/h"),
			loader.LoadDynamicSprite("item", "2/attack/v")},
	}
	sprites[4][1].(*sprite.DynamicSprite).Once = true //攻击只播放一次
	res.WalkEnemy = NewWalkEnemy(res.triggerAttack, sprites)
	factory.FillPosObject(res.PosObject, o)
	res.typeName = "enemy2"
	res.speed = 1
	res.actions = append(res.actions, res.attack)
	res.dieSpriteNames = []string{"2/die/move", "2/die/idle"}
	return res
}

// Enemy2 喷火怪
type Enemy2 struct {
	*WalkEnemy
	attackTimer int // 攻击动画持续时间
}

func (e *Enemy2) triggerAttack() {
	if player == nil {
		e.randState() // 若玩家不为nil 根据位置判断技能方向
	} else if e.Y-player.Y > math.Abs(player.X-e.X) {
		e.switchState(4, 1)
		e.attackTimer = 15
		tool.AddObject("enemy", NewFire(0, e.X, e.Y-16))
	} else {
		e.ScaleX = utils.Sign(player.X - e.X)
		e.switchState(4, 0)
		e.attackTimer = 15
		tool.AddObject("enemy", NewFire(int(e.ScaleX), e.X, e.Y-16))
	}
}

func (e *Enemy2) attack() {
	e.attackTimer--
	if e.attackTimer < 0 {
		e.randState()
	}
}
