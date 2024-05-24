package other

import (
	"gameBase/loader"
	"gameBase/model"
	"gameBase/object"
	"gameBase/other"
	"gameBase/sprite"
	"gameBase/tool"
)

type PlayerDie struct {
	*object.TimerEffect
}

func (p *PlayerDie) Update() {
	p.Y--
}

func (p *PlayerDie) animEnd() {
	p.Sprite.(*sprite.DynamicSprite).SetIndex(1)
}

func (p *PlayerDie) effectEnd() {
	tool.AddObject("player", NewPlayerBorn())
}

func NewPlayerDie(x, y float64) *PlayerDie {
	res := new(PlayerDie)
	res.TimerEffect = object.NewTimerEffect(60, x, y)
	spr := loader.LoadDynamicSprite("item", "payer/die")
	spr.AnimEnd = res.animEnd
	res.Sprite = spr
	res.EffectEnd = res.effectEnd
	return res
}

type PlayerBorn struct {
	*object.TimerEffect
	snowSprites []other.Sprite
}

func NewPlayerBorn() *PlayerBorn {
	res := new(PlayerBorn)
	res.TimerEffect = object.NewTimerEffect(60, 64, 192)
	res.snowSprites = []other.Sprite{loader.LoadDynamicSprite("item", "payer/born/effect1"),
		loader.LoadDynamicSprite("item", "payer/born/effect2")}
	res.Sprite = res.snowSprites[0]
	res.EffectEnd = res.effectEnd
	return res
}

func (p *PlayerBorn) Update() {
	if p.Timer == 30 {
		p.Sprite = p.snowSprites[1]
	}
}

func (p *PlayerBorn) effectEnd() {
	pos := model.NewPointObject()
	pos.X, pos.Y = p.X, p.Y
	pos.Visible = true
	tool.AddObject("player", createPlayer(pos, nil))
}

type FireBoom struct {
	*object.TimerEffect
}

func NewFireBoom(x, y float64) *FireBoom {
	res := new(FireBoom)
	res.TimerEffect = object.NewTimerEffect(15, x, y)
	res.Sprite = loader.LoadStaticSprite("item", "2/bullet/boom")
	return res
}

type DieEnemy struct {
	*object.CollisionObject
	xSpeed, ySpeed float64
	idleTimer      int
	idle           bool
	die            bool
	idleSprite     other.Sprite
}

// CanCollision 禁用其他人都它的碰撞检测
func (d *DieEnemy) CanCollision(tag int) bool {
	return false
}

func (d *DieEnemy) IsDie() bool {
	return d.die
}

func (d *DieEnemy) Update() {
	if d.idle {
		d.idleTimer--
		if d.idleTimer < 0 {
			d.die = true
			tool.AddObject("player", NewPotion(d.X, d.Y))
		}
	} else {
		if d.MoveHorizontal("collision", WallTag, d.xSpeed) {
			d.xSpeed *= -1
		}
		if d.ySpeed < MaxSpeed {
			d.ySpeed += G
		}
		if d.ySpeed > 0 {
			if d.MoveVertical("collision", GroundTag|FloorTag, d.ySpeed) {
				d.idle = true
				d.Sprite = d.idleSprite
			}
		} else {
			d.Y += d.ySpeed
		}
	}
}

func NewDieEnemy(x, y, dir float64, move, idle other.Sprite) *DieEnemy {
	res := new(DieEnemy)
	res.CollisionObject = object.NewCollisionObject()
	res.idle = false
	res.idleTimer = 30
	res.xSpeed, res.ySpeed = dir*2, -3
	res.X, res.Y = x, y
	res.die = false
	res.Sprite = move
	res.idleSprite = idle
	res.OffsetX, res.OffsetY = -6, -6
	res.Width, res.Height = 12, 6
	return res
}
