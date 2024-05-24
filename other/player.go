package other

import (
	"gameBase/config"
	"gameBase/factory"
	"gameBase/loader"
	"gameBase/model"
	"gameBase/object"
	"gameBase/other"
	"gameBase/sprite"
	"gameBase/tool"
	"gameBase/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	factory.RegisterPosFactory("player", createPlayer)
}

var (
	player *Player // 向外暴露
)

func dynamicSprite(action string) other.Sprite {
	return loader.LoadDynamicSprite("item", "payer/"+action)
}

func staticSprite(action string) other.Sprite {
	return loader.LoadStaticSprite("item", "payer/"+action)
}

func createPlayer(o *model.PointObject, m *model.Map) other.GameObject {
	res := new(Player)
	res.CollisionObject = object.NewCollisionObject()
	factory.FillPosObject(res.PosObject, o)
	res.Width, res.Height = 10, 12
	res.OffsetX, res.OffsetY = -5, -12
	res.sprites = [][]other.Sprite{
		{staticSprite("idle/nomal"), dynamicSprite("idle/fast")},
		{dynamicSprite("move/nomal"), dynamicSprite("move/fast")},
		{dynamicSprite("shoot"), staticSprite("kick")},
		{dynamicSprite("jump/start"), staticSprite("jump/loop")},
		{dynamicSprite("big")},
		{dynamicSprite("move/push")},
		{staticSprite("idle/nomal")},
	}
	res.actions = []func(){res.idle, res.move, res.attack, res.jump, res.skill, res.push, res.follow}
	res.keys = []ebiten.Key{ebiten.KeyA, ebiten.KeyD, ebiten.KeyJ, ebiten.KeyK}
	res.speeds = []float64{1, 2}
	res.Tag = PlayerTag
	res.fast = 0
	res.attackSizeType = 0
	res.attackTimeType = 0
	res.die = false
	res.protectTimer = 60
	afterPropertySet(res)
	player = res
	return res
}

func afterPropertySet(res *Player) {
	res.switchState(0, 0)
	res.sprites[4][0].(*sprite.DynamicSprite).AnimEnd = res.animEnd
	res.sprites[3][0].(*sprite.DynamicSprite).AnimEnd = res.animEnd
	res.sprites[2][0].(*sprite.DynamicSprite).AnimEnd = res.animEnd
}

type Player struct {
	*object.CollisionObject
	// idle  normal  fast
	// move  normal  fast
	// attack normal kick
	// jump   start loop
	// skill
	// push
	// follow
	sprites        [][]other.Sprite
	state, typ     int
	actions        []func()
	dir            int
	keys           []ebiten.Key
	fast           int       // 是否快速
	speeds         []float64 // 两种速度
	ySpeed         float64
	attackSizeType int
	attackTimeType int
	attackTimer    int
	die            bool
	skillTimer     int
	protectTimer   int
}

func (p *Player) Priority() int {
	return -2233
}

func (p *Player) IsDie() bool {
	return p.die
}

func (p *Player) Update() {
	if utils.IsKeyPress(p.keys[LeftKey]) {
		p.dir = -1
	} else if utils.IsKeyPress(p.keys[RightKey]) {
		p.dir = 1
	} else {
		p.dir = 0
	}
	if p.dir != 0 {
		p.ScaleX = float64(p.dir)
	}
	p.attackTimer--
	p.protectTimer--
	p.actions[p.state]()
}

func (p *Player) switchState(state int, typ int) {
	p.state, p.typ = state, typ
	p.Sprite = p.sprites[state][typ]
	if spr, ok := p.Sprite.(*sprite.DynamicSprite); ok {
		spr.Reset()
	}
}

func (p *Player) idle() {
	p.checkAir()
	p.checkMove()
	p.checkJump()
	p.checkAttack()
}

func (p *Player) move() {
	p.checkAir()
	p.MoveHorizontal("collision", WallTag, p.ScaleX*p.speeds[p.fast])
	p.checkPush() // 只有移动状态可能转换为 推状态
	p.checkIdle()
	p.checkJump()
	p.checkAttack()
}

func (p *Player) attack() {
	if p.typ == 1 && p.attackTimer < 0 {
		p.switchState(0, p.fast)
	}
}

func (p *Player) jump() {
	p.checkAttack()
	p.MoveHorizontal("collision", WallTag, float64(p.dir)*p.speeds[p.fast])
	if p.ySpeed < MaxSpeed {
		p.ySpeed += G
	}
	if p.ySpeed > 0 {
		if p.checkGround(p.ySpeed) {
			p.Y = float64(int(p.Y))
			for !p.checkGround(1) {
				p.Y++
			}
			p.switchState(0, p.fast)
		} else {
			p.Y += p.ySpeed
		}
	} else {
		p.Y += p.ySpeed
	}
}

var (
	moveKeys = []ebiten.Key{ebiten.KeyW, ebiten.KeyS, ebiten.KeyA, ebiten.KeyD}
	dirs     = [][]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
)

func (p *Player) skill() {
	p.skillTimer--
	if p.skillTimer < 0 {
		p.switchState(0, p.fast)
	}
	for i := 0; i < len(moveKeys); i++ {
		if utils.IsKeyPress(moveKeys[i]) {
			p.X += float64(dirs[i][0]) * 1.5
			p.Y += float64(dirs[i][1]) * 1.5
			break
		}
	}
	p.X = utils.Clamp(p.X, 6, float64(config.Width-6))
	p.Y = utils.Clamp(p.Y, 12, float64(config.Height-33))
	if bump, ok := tool.CollisionRect("enemy", EnemyTag, p).(BumpAble); ok {
		bump.Bump(p.ScaleX)
	}
}

func (p *Player) checkAir() {
	if !p.checkGround(1) {
		p.ySpeed = 0
		p.switchState(3, 1)
	}
}

func (p *Player) checkMove() {
	if p.dir != 0 {
		p.switchState(1, p.fast)
	}
}

func (p *Player) checkIdle() {
	if p.dir == 0 {
		p.switchState(0, p.fast)
	}
}

func (p *Player) checkJump() {
	if utils.IsKeyPress(p.keys[JumpKey]) {
		p.ySpeed = JumpSpeed
		p.switchState(3, 0)
	}
}

// Divide 对外 触发跳跃方法
func (p *Player) Divide() {
	p.ySpeed = JumpSpeed
	p.protectTimer = 60
	p.switchState(3, 0)
}

func (p *Player) animEnd() {
	switch p.state {
	case 3:
		p.switchState(3, 1)
	case 2:
		p.switchState(0, p.fast)
	case 4:
		p.Sprite.(*sprite.DynamicSprite).SetIndex(1)
	}
}

func (p *Player) checkGround(offset float64) bool {
	if tool.CollisionPoint("collision", GroundTag|FloorTag|SnowBallTag, p.X, p.GetBottom()) != nil {
		return false //若自身在墙里  不能算在地面上
	} // 墙 雪球都可以站立
	return p.CollisionBottom("collision", GroundTag|FloorTag|SnowBallTag, offset) != nil
}

func (p *Player) checkAttack() {
	if p.attackTimer < 0 && utils.IsKeyPress(p.keys[AttackKey]) {
		tool.AddObject("player", NewSnow(p.X, p.Y-13, p.attackSizeType, p.attackTimeType, p.ScaleX))
		p.switchState(2, 0)
		p.attackTimer = 15
	}
}

// Die 玩家死亡逻辑
func (p *Player) Die() {
	if p.state == 6 || p.state == 4 || p.protectTimer > 0 {
		return // 跟随模式技能时间 是无敌的
	}
	tool.AddObject("player", NewPlayerDie(p.X, p.Y-12))
	player = nil
	p.die = true
}

func (p *Player) checkKick(snowBall *SnowBall) {
	if utils.IsKeyPress(p.keys[AttackKey]) {
		snowBall.kick(p.ScaleX)
		p.switchState(2, 1)
		p.attackTimer = 5
	}
}

func (p *Player) push() {
	p.checkIdle()
	if snowBall, ok :=
		tool.CollisionPoint("collision", SnowBallTag, p.X+p.ScaleX*6, p.Y-6).(*SnowBall); ok {
		p.checkKick(snowBall)
		if snowBall.push(p.ScaleX) { // 推雪球时雪球速度 0.5
			p.X += 0.5 * p.ScaleX
		}
	} else {
		p.switchState(0, p.fast)
	}
}

func (p *Player) checkPush() {
	if tool.CollisionPoint("collision", SnowBallTag, p.X+p.ScaleX*6, p.Y-6) != nil {
		p.switchState(5, 0)
	}
}

// 跟随雪球模式  由雪球控制移动
func (p *Player) follow() {

}

func (p *Player) GetPotion(typ int) {
	switch typ {
	case 0:
		p.attackSizeType = 1
	case 1:
		p.fast = 1
	case 2:
		p.attackTimeType = 1
	case 3:
		p.switchState(4, 0)
		p.skillTimer = 360
	default:
		utils.PanicF("未知的药水类型:%d", typ)
	}
}
