package other

import (
	"gameBase/factory"
	"gameBase/loader"
	"gameBase/model"
	"gameBase/other"
	"gameBase/tool"
	"gameBase/utils"
	"math/rand"
)

func init() {
	factory.RegisterPosFactory("boss", createBoss)
}

func createBoss(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Boss)
	res.Enemy = NewEnemy()
	factory.FillPosObject(res.PosObject, o)
	die := loader.LoadDynamicSprite("item", "boss/die")
	die.AnimEnd = res.animEnd
	res.sprites = []other.Sprite{loader.LoadStaticSprite("item", "boss/idle"),
		loader.LoadStaticSprite("item", "boss/jump"), die}
	res.actions = []func(){res.idle, res.jump, utils.EmptyFunc}
	res.hp = 10800
	res.Tag = BossTag
	res.OffsetX, res.OffsetY = 0, -60
	res.Width, res.Height = 28, 60
	res.idleTimer = utils.RandInt(30, 90)
	res.switchState(0)
	res.heights = []float64{112, 192}
	return res
}

type Boss struct {
	*Enemy
	state        int // idle jump die
	sprites      []other.Sprite
	actions      []func()
	hp           int
	ySpeed       float64
	idleTimer    int
	heights      []float64
	targetHeight float64
	lastYSpeed   float64
}

// Hurt 雪球  1800  攻击按普通的    默认6个雪球解决
func (e *Boss) Hurt(value int) {
	e.hp -= value
	if e.hp < 0 {
		e.switchState(2)
	}
}

func (e *Boss) Update() {
	e.actions[e.state]()
	e.Enemy.Update()
}

func (e *Boss) switchState(state int) {
	e.state = state
	e.Sprite = e.sprites[state]
}

func (e *Boss) animEnd() {
	e.die = true // 只有死亡有动画
}

func (e *Boss) idle() {
	e.idleTimer--
	if e.idleTimer < 0 {
		e.targetHeight = e.heights[rand.Intn(2)]
		e.ySpeed = -4
		if e.targetHeight < e.Y-32 {
			e.ySpeed -= 3
		}
		e.switchState(1)
	}
}

func (e *Boss) jump() {
	if e.ySpeed < MaxSpeed {
		e.lastYSpeed = e.ySpeed
		e.ySpeed += G // 最高点尝试仍球
		if e.lastYSpeed*e.ySpeed < 0 && rand.Intn(3) == 0 {
			tool.AddObject("enemy", NewEgg(e.X+8, e.Y-40))
		}
	}
	if e.ySpeed > 0 { // 下落时才判断高度
		if e.Y+e.ySpeed > e.targetHeight {
			e.Y = e.targetHeight
			e.idleTimer = rand.Intn(60) + 30
			e.switchState(0)
		} else {
			e.Y += e.ySpeed
		}
	} else {
		e.Y += e.ySpeed
	}
}

var (
	enemyNames = []string{"enemy1", "enemy2", "enemy3", "enemy4", "enemy5", "enemy6", "son"}
)

type Egg struct {
	*Enemy
	xSpeed, ySpeed float64
	bronTimer      int
	fly            bool
	bronSprite     other.Sprite
}

func (e *Egg) Hurt(value int) {
	e.die = true // 生成雪球  蛋被攻击会失去随机性变为最基本的蛋
	ball := NewSnowBall(e.X, e.Y, "son")
	ball.Name = tool.AddObject("enemy", ball)
}

// Bump 蛋状态下没有奖励
func (e *Egg) Bump(dir float64) {
	e.die = true
}

func (e *Egg) Update() {
	e.Enemy.Update()
	if e.fly {
		e.X += e.xSpeed
		if e.ySpeed < MaxSpeed {
			e.ySpeed += G
		}
		if e.ySpeed > 0 {
			if e.MoveVertical("collision", FloorTag|GroundTag, e.ySpeed) {
				e.fly = false
			}
		} else {
			e.Y += e.ySpeed
		}
	} else {
		e.bronTimer--
		if e.bronTimer < 0 {
			e.die = true
			enemy := createEnemy(enemyNames[rand.Intn(len(enemyNames))], e.X, e.Y)
			tool.AddObject("enemy", enemy)
		}
	}
}

func NewEgg(x, y float64) *Egg {
	res := new(Egg)
	res.Enemy = NewEnemy()
	res.X, res.Y = x, y
	res.xSpeed, res.ySpeed = -utils.RandFloat(3, 6), -utils.RandFloat(1, 2)
	res.bronTimer = 30
	res.fly = true
	res.Sprite = loader.LoadStaticSprite("item", "boss/son/egg")
	res.bronSprite = loader.LoadStaticSprite("item", "boss/son/bron")
	res.OffsetX, res.OffsetY = -8, -16
	res.Width, res.Height = 16, 16
	return res
}

type Son struct {
	*Enemy
	ySpeed float64
}

func (s *Son) Hurt(value int) {
	s.die = true
	ball := NewSnowBall(s.X, s.Y, "son")
	ball.Name = tool.AddObject("enemy", ball)
}

func (s *Son) Bump(dir float64) {
	s.die = true
	tool.AddObject("player", NewPotion(s.X, s.Y))
}

func (s *Son) Update() {
	s.Enemy.Update()
	if s.ySpeed < MaxSpeed {
		s.ySpeed += G
	}
	if s.MoveVertical("collision", GroundTag|FloorTag, s.ySpeed) {
		s.ySpeed = 0
	}
	if s.MoveHorizontal("collision", WallTag, s.ScaleX) {
		s.ScaleX *= -1
	}
}

func NewSon(x, y float64) *Son {
	res := new(Son)
	res.Enemy = NewEnemy()
	res.X, res.Y = x, y
	res.Sprite = loader.LoadDynamicSprite("item", "boss/son/move")
	res.OffsetX, res.OffsetY = -8, -16
	res.Width, res.Height = 16, 16
	return res
}
