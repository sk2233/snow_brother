package other

import (
	"gameBase/loader"
	"gameBase/model"
	"gameBase/object"
	"gameBase/other"
	"gameBase/sprite"
	"gameBase/tool"
	"gameBase/utils"
	"math/rand"
)

// Enemy 敌人基类 都是受玩家攻击影响的对象
type Enemy struct {
	*object.CollisionObject
	die bool
}

func (e *Enemy) IsDie() bool {
	return e.die
}

func (e *Enemy) Hurt(value int) {

}

func (e *Enemy) Update() {
	if player != nil && player.CollisionRect(e) {
		player.Die()
	}
}

func NewEnemy() *Enemy {
	res := new(Enemy)
	res.CollisionObject = object.NewCollisionObject()
	res.Tag = EnemyTag
	return res
}

func createEnemy(typ string, x, y float64) other.GameObject {
	pos := model.NewPointObject()
	pos.X, pos.Y = x, y
	pos.Visible = true
	switch typ {
	case "enemy1":
		return createEnemy1(pos, nil)
	case "enemy2":
		return createEnemy2(pos, nil)
	case "enemy3":
		return createEnemy3(pos, nil)
	case "enemy4":
		return createEnemy4(pos, nil)
	case "enemy5":
		return createEnemy5(pos, nil)
	case "enemy6":
		return createEnemy6(pos, nil)
	case "son":
		return NewSon(x, y)
	default:
		utils.PanicF("未知的类型%s", typ)
		return nil
	}
}

// WalkEnemy 最常见
type WalkEnemy struct {
	*Enemy
	actions []func()         // 子类指定
	sprites [][]other.Sprite // 为了兼容部分图片有两个的  后代初始化
	// idle move jump shake  other
	dieSpriteNames []string // 为2的数组 先 move 再idle
	state          int
	typeName       string // 需要赋值
	timer          int    // idle move 等平静时间的时间  每 4-8s触发一次
	skillStart     func()
	ySpeed         float64
	shakeTimer     int
	speed          float64
}

// Bump 默认会被撞死
func (w *WalkEnemy) Bump(dir float64) { // 延迟加载
	move := loader.LoadDynamicSprite("item", w.dieSpriteNames[0])
	idle := loader.LoadStaticSprite("item", w.dieSpriteNames[1])
	tool.AddObject("player", NewDieEnemy(w.X, w.Y, dir, move, idle))
	w.die = true
}

func (w *WalkEnemy) Shake() {
	w.shakeTimer = 60
	w.switchState(3, 0)
}

func NewWalkEnemy(skillStart func(), sprites [][]other.Sprite) *WalkEnemy {
	res := new(WalkEnemy)
	res.Enemy = NewEnemy() // 行走的敌人 默认一样大
	res.OffsetX, res.OffsetY = -5, -14
	res.Width, res.Height = 10, 14
	res.skillStart = skillStart
	res.sprites = sprites
	res.actions = []func(){res.idle, res.move, res.jump, res.shake}
	//res.resetTimer() // 先赋值 再调用方法
	res.switchState(0, 0)
	return res
}

func (w *WalkEnemy) Hurt(value int) {
	w.die = true // 生成雪球
	ball := NewSnowBall(w.X, w.Y, w.typeName)
	ball.Name = tool.AddObject("enemy", ball)
}

func (w *WalkEnemy) resetTimer() {
	w.timer = rand.Intn(4*60) + 4*60
}

func (w *WalkEnemy) Update() {
	w.Enemy.Update()
	w.actions[w.state]()
}

// typ 用于 哪些有多张图片的
func (w *WalkEnemy) switchState(state, typ int) {
	w.state = state
	w.Sprite = w.sprites[state][typ]
	if dynamic, ok := w.Sprite.(*sprite.DynamicSprite); ok {
		dynamic.Reset()
	}
}

func (w *WalkEnemy) idle() {
	w.passTime()
	w.checkAir()
}

func (w *WalkEnemy) move() {
	w.passTime()
	if w.MoveHorizontal("collision", WallTag, w.ScaleX*w.speed) {
		w.ScaleX *= -1
	}
	w.checkAir()
}

func (w *WalkEnemy) jump() {
	if w.ySpeed < MaxSpeed {
		w.ySpeed += G
	}
	if w.ySpeed > 0 {
		if w.checkGround(w.ySpeed) {
			w.Y = float64(int(w.Y))
			for !w.checkGround(1) {
				w.Y++
			}
			w.randState()
		} else {
			w.Y += w.ySpeed
		}
	} else {
		w.Y += w.ySpeed
	}
}

func (w *WalkEnemy) shake() {
	w.shakeTimer--
	if w.shakeTimer < 0 {
		w.randState()
	}
	w.checkAir()
}

func (w *WalkEnemy) passTime() {
	w.timer--
	if w.timer < 0 {
		w.resetTimer()
		w.randState()
	}
}

func (w *WalkEnemy) randState() {
	if tool.CollisionPoint("collision", GroundTag|FloorTag, w.X, w.Y-28) != nil {
		w.ySpeed = JumpSpeed
		w.switchState(2, 0)
	} else { // 能跳的话  优先跳跃
		switch rand.Intn(4) {
		case 0:
			w.timer = 30 //最多站0.5s
			w.switchState(0, 0)
		case 1:
			w.ScaleX = 1
			w.switchState(1, 0)
		case 2:
			w.ScaleX = -1
			w.switchState(1, 0)
		case 3:
			w.skillStart()
		}
	}
}

func (w *WalkEnemy) checkGround(offset float64) bool {
	if tool.CollisionPoint("collision", GroundTag|FloorTag, w.X, w.GetBottom()) != nil {
		return false //若自身在墙里  不能算在地面上
	} // 墙 雪球都可以站立
	return w.CollisionBottom("collision", GroundTag|FloorTag, offset) != nil
}

func (w *WalkEnemy) checkAir() {
	if !w.checkGround(1) {
		w.ySpeed = 0 // 有down图片的 使用down
		w.switchState(2, len(w.sprites[2])-1)
	}
}
