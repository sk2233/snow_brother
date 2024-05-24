package other

import (
	"gameBase/loader"
	"gameBase/object"
	"gameBase/other"
	"gameBase/tool"
	"gameBase/utils"
	"strconv"
)

var (
	powers []int
	timers []int
	paths  []string
)

func init() {
	powers = []int{150, 300}
	timers = []int{60, 90}
	paths = []string{"payer/bullet/nomal", "payer/bullet/big"}
}

type Snow struct {
	*object.CollisionObject
	Power  int //威力大小
	Timer  int // 飞行时间
	Die    bool
	ySpeed float64
}

func (s *Snow) IsDie() bool {
	return s.Die
}

func (s *Snow) Update() {
	if s.Timer > 10 {
		s.X += s.ScaleX * 2
	} else if s.Timer > 0 {
		s.X += s.ScaleX * 2
		s.ySpeed += G
		s.Y += s.ySpeed
	} else {
		s.Die = true
	}
	s.Timer--
	if hurt, ok := tool.CollisionRect("enemy", EnemyTag|SnowBallTag|BossTag, s).(HurtAble); ok {
		hurt.Hurt(s.Power)
		s.Die = true
	}
}

// NewSnow sizeType 0 小 1 大  timeType 0 正常  1 长
func NewSnow(x, y float64, sizeType, timeType int, dir float64) *Snow {
	res := new(Snow)
	res.CollisionObject = object.NewCollisionObject()
	res.X, res.Y = x, y
	res.Tag = SnowTag
	res.ScaleX = dir // 直接记录方向
	res.Power = powers[sizeType]
	res.Timer = timers[timeType]
	res.Die = false
	res.Sprite = loader.LoadDynamicSprite("item", paths[sizeType])
	res.OffsetX, res.OffsetY = -float64(sizeType+1), -float64(sizeType+1)*2
	res.Width, res.Height = float64(sizeType+1)*2, float64(sizeType+1)*4
	return res
}

type SnowBall struct {
	*object.CollisionObject
	enemyType    string // 敌人类型融合时恢复
	move         bool
	moveSprite   other.Sprite
	meltSprites  []other.Sprite
	timer        int
	die          bool
	Name         string
	ySpeed       float64
	player       *Player
	collisionNum int
}

func (s *SnowBall) Hurt(value int) {
	s.timer += value // 因为图层移动 不会承受过多
	s.Sprite = s.meltSprites[utils.Min(2, s.timer/300)]
	if s.timer > 600 {
		tool.MoveObjectLayer("enemy", "collision", s.Name)
	}
}

func (s *SnowBall) kick(dir float64) {
	s.ScaleX = dir
	s.Sprite = s.moveSprite
	tool.MoveObjectLayer("collision", "player", s.Name)
	s.move = true
}

func (s *SnowBall) IsDie() bool {
	return s.die
}

// 0  300   600

func (s *SnowBall) Update() {
	// 降落判断
	if s.CollisionBottom("collision", GroundTag|FloorTag, 1) == nil {
		if s.ySpeed < MaxSpeed {
			s.ySpeed += G * 2
		}
		if s.MoveVertical("collision", GroundTag|FloorTag, s.ySpeed) {
			s.ySpeed = 0
		}
	}
	if s.move { //滚动过程
		s.handlePlayer()
		s.handleMove()
		s.handleSnowBall()
		s.handleEnemy()
		s.handleBallExit()
	} else { // 融化过程
		if s.timer%300 == 0 {
			if s.timer == 0 {
				s.die = true //恢复
				enemy := createEnemy(s.enemyType, s.X, s.Y)
				if shake, ok := enemy.(ShakeAble); ok {
					shake.Shake() // boss的孩子不会抖雪
				}
				tool.AddObject("enemy", enemy)
			} else {
				index := utils.Min(s.timer/300-1, 2)
				s.Sprite = s.meltSprites[index]
				if s.timer == 600 { // 移动图层
					tool.MoveObjectLayer("collision", "enemy", s.Name)
				}
			}
		}
		s.timer--
	}
}

// 玩家推雪球 也是 0.5
func (s *SnowBall) push(dir float64) bool {
	if tool.CollisionPoint("collision", WallTag|SnowBallTag, s.X+dir*6.5, s.Y-9) != nil {
		return false
	}
	s.X += dir * 0.5
	return true
}

func (s *SnowBall) handlePlayer() {
	if player != nil && player.CollisionRect(s) {
		player.switchState(6, 0)
		s.player = player
	}
	if s.player != nil {
		s.player.Y = s.Y
		s.player.X = s.X - s.ScaleX*5
		s.player.ScaleX = -s.ScaleX
	}
}

func (s *SnowBall) handleMove() {
	if s.MoveHorizontal("collision", WallTag, s.ScaleX*3) {
		s.collisionNum-- // 最大碰撞次数处理
		if s.collisionNum < 0 {
			s.death()
		} else {
			s.ScaleX *= -1
		}
	}
}

func (s *SnowBall) handleEnemy() {
	if bump, ok := tool.CollisionPoint("enemy", EnemyTag, s.X, s.Y).(BumpAble); ok {
		bump.Bump(s.ScaleX)
	}
	if hurt, ok := tool.CollisionPoint("enemy", BossTag, s.X, s.Y).(HurtAble); ok {
		hurt.Hurt(1800)
		s.death()
	}
}

func (s *SnowBall) handleBallExit() {
	if tool.CollisionPoint("collision", BallExitTag, s.X, s.Y-16) != nil {
		s.death()
	}
}

func (s *SnowBall) death() {
	s.die = true
	if s.player != nil {
		s.player.Divide()
	}
}

func (s *SnowBall) handleSnowBall() {
	if snowBall, ok :=
		tool.CollisionPoint("collision", SnowBallTag, s.X-s.ScaleX*6, s.Y-16).(*SnowBall); ok {
		snowBall.kick(-s.ScaleX)
	}
}

// NewSnowBall 默认生成在敌人层  圆满时移动到碰撞层
func NewSnowBall(x, y float64, enemyType string) *SnowBall {
	res := new(SnowBall)
	res.CollisionObject = object.NewCollisionObject()
	res.X, res.Y = x, y
	res.enemyType = enemyType
	res.Tag = SnowBallTag
	res.move = false
	res.die = false
	res.moveSprite = loader.LoadDynamicSprite("item", "ball/roll")
	res.meltSprites = make([]other.Sprite, 3)
	for i := 0; i < 3; i++ {
		res.meltSprites[i] = loader.LoadStaticSprite("item", "ball/melt/state"+strconv.Itoa(i+1))
	}
	res.OffsetX, res.OffsetY = -6, -18
	res.Width, res.Height = 12, 18
	res.Hurt(300)
	res.ySpeed = 0
	res.collisionNum = 6
	return res
}
