package other

import (
	"gameBase/loader"
	"gameBase/object"
	"gameBase/tool"
)

// 敌人的子弹  集  不对玩家攻击响应

type Bullet struct {
	*object.CollisionObject
	die    bool
	dieFun func()
}

func NewBullet(x, y float64) *Bullet {
	res := new(Bullet)
	res.CollisionObject = object.NewCollisionObject()
	res.Tag = BulletTag
	res.X, res.Y = x, y
	res.die = false
	res.ExitRoom = res.exitRoom
	return res
}

func (e *Bullet) IsDie() bool {
	return e.die
}

func (e *Bullet) Update() {
	if player != nil && player.CollisionRect(e) {
		player.Die()
		e.die = true
		if e.dieFun != nil {
			e.dieFun()
		}
	}
}

// 退出房间自动销毁
func (e *Bullet) exitRoom() {
	e.die = true
}

type Fire struct {
	*Bullet
	xSpeed, ySpeed float64
}

func (e *Fire) Update() {
	e.X += e.xSpeed
	e.Y += e.ySpeed
	e.Bullet.Update()
}

func (e *Fire) createBoom() {
	tool.AddObject("enemy", NewFireBoom(e.X, e.Y))
}

// NewFire dir -1 1 横向   0 向下
func NewFire(dir int, x, y float64) *Fire {
	res := new(Fire)
	res.Bullet = NewBullet(x, y)
	if dir == 0 {
		res.xSpeed = 0
		res.ySpeed = 2
		res.Sprite = loader.LoadDynamicSprite("item", "2/bullet/v")
	} else {
		res.ScaleX = float64(dir)
		res.xSpeed = res.ScaleX * 2
		res.ySpeed = 0
		res.Sprite = loader.LoadStaticSprite("item", "2/bullet/h")
	}
	res.OffsetX, res.OffsetY = -6, -6
	res.Width, res.Height = 12, 12
	res.dieFun = res.createBoom
	return res
}

type Knife struct {
	*Bullet
	xSpeed float64
}

func (e *Knife) Update() {
	e.X += e.xSpeed
	e.Bullet.Update()
}

// NewKnife dir -1 1 横向
func NewKnife(dir float64, x, y float64) *Knife {
	res := new(Knife)
	res.Bullet = NewBullet(x, y)
	res.ScaleX = dir
	res.xSpeed = dir * 2
	res.Sprite = loader.LoadDynamicSprite("item", "4/bullet")
	res.OffsetX, res.OffsetY = -4, -1
	res.Width, res.Height = 8, 2
	return res
}
