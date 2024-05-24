package other

// 可以考虑使用预组合常量  方便后面添加新标签
const ( // 推荐使用常量标识对应质数
	GroundTag = 2 << iota // 不可穿透的地面
	FloorTag              // 可以穿透的地面
	WallTag               // 水平墙壁
	PlayerTag
	SnowBallTag
	SnowTag
	EnemyTag
	BulletTag   // 敌人子弹
	BallExitTag // 雪球消失点
	PotionTag   // 药水
	BossTag
)

const (
	LeftKey = iota
	RightKey
	AttackKey
	JumpKey
)

const (
	G         = 0.2
	MaxSpeed  = 4.0  // 防止穿透检测
	JumpSpeed = -4.0 // 必须是负数 向上跳
)

// HurtAble 所有对主角攻击有效的对象接口
type HurtAble interface {
	Hurt(int)
}

// ShakeAble 可以抖雪的对象
type ShakeAble interface {
	Shake()
}

// BumpAble 可被雪球撞击的
type BumpAble interface {
	Bump(dir float64) //传入撞击方向
}
