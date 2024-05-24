package other

import (
	"gameBase/factory"
	"gameBase/loader"
	"gameBase/model"
	"gameBase/other"
	"gameBase/tool"
	"gameBase/utils"
)

func init() {
	factory.RegisterPosFactory("enemy5", createEnemy5)
}

func createEnemy5(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Enemy5)
	sprites := [][]other.Sprite{
		{loader.LoadStaticSprite("item", "5/idle")},
		{loader.LoadDynamicSprite("item", "5/move")},
		{loader.LoadStaticSprite("item", "5/jump")},
		{loader.LoadDynamicSprite("item", "5/shake")},
		{loader.LoadStaticSprite("item", "5/down")}, // 穿墙
		{loader.LoadStaticSprite("item", "5/fire")}, // 放火
	}
	res.WalkEnemy = NewWalkEnemy(res.trigger, sprites)
	factory.FillPosObject(res.PosObject, o)
	res.typeName = "enemy5"
	res.speed = 1
	res.actions = append(res.actions, res.action, res.action)
	res.dieSpriteNames = []string{"5/die/move", "5/die/idle"}
	return res
}

// Enemy5 猴子
type Enemy5 struct {
	*WalkEnemy
	actionTimer int // 吊着或放火动画持续时间
}

func (e *Enemy5) trigger() { // 必须当前允许穿墙才能穿墙
	if tool.CollisionPoint("collision", FloorTag, e.X, e.GetBottom()+1) != nil {
		e.Y += 26
		e.switchState(4, 0)
	} else {
		if player != nil { // 火是持续的 一定要放
			e.ScaleX = utils.Sign(player.X - e.X)
		}
		tool.AddObject("enemy", NewBonfire(e.X, e.Y-14, e.ScaleX))
		e.switchState(5, 0)
	}
	e.actionTimer = 30
}

// 都是统一时间过了就恢复原来状态
func (e *Enemy5) action() {
	e.actionTimer--
	if e.actionTimer < 0 {
		e.switchState(0, 0)
	}
}

type Bonfire struct {
	*Enemy
	timer          int //10秒
	ground         bool
	xSpeed, ySpeed float64
	fireSprite     other.Sprite
}

func (e *Bonfire) Hurt(value int) {
	e.timer -= value / 3
}

func (e *Bonfire) Update() {
	if e.ground {
		if player != nil && player.CollisionRect(e) {
			player.Die()
		}
		e.timer--
		if e.timer < 0 {
			e.die = true
		}
	} else {
		if e.ySpeed < MaxSpeed {
			e.ySpeed += G
		}
		if e.ySpeed > 0 {
			if e.MoveVertical("collision", GroundTag|FloorTag, e.ySpeed) {
				e.ground = true
				e.Sprite = e.fireSprite
			}
		} else {
			e.Y += e.ySpeed
		}
		e.X += e.xSpeed
	}
}

func NewBonfire(x, y, dir float64) *Bonfire {
	res := new(Bonfire)
	res.Enemy = NewEnemy()
	res.timer = 600
	res.ground = false
	res.X, res.Y = x, y
	res.Sprite = loader.LoadStaticSprite("item", "5/bullet/seed")
	res.fireSprite = loader.LoadDynamicSprite("item", "5/bullet/fire")
	res.xSpeed = dir * 2
	res.ySpeed = -3
	res.OffsetX, res.OffsetY = -6, -18
	res.Width, res.Height = 12, 18
	return res
}
