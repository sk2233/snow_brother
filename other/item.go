package other

import (
	"gameBase/loader"
	"gameBase/object"
	"math/rand"
)

type Potion struct {
	*object.CollisionObject
	Type  int
	timer int
	die   bool
}

func (p *Potion) Update() {
	p.timer--
	if p.timer < 0 {
		p.die = true
	}
	if player != nil && p.CollisionRect(player) {
		p.die = true
		player.GetPotion(p.Type)
	}
}

func (p *Potion) IsDie() bool {
	return p.die
}

var (
	spritePaths = []string{"item/blue", "item/red", "item/yellow", "item/green"}
)

func NewPotion(x, y float64) *Potion {
	res := new(Potion)
	res.CollisionObject = object.NewCollisionObject()
	res.Type = rand.Intn(len(spritePaths))
	res.X, res.Y = x, y
	res.Sprite = loader.LoadDynamicSprite("item", spritePaths[res.Type])
	res.Tag = PotionTag
	res.OffsetX, res.OffsetY = -7, -9
	res.Width, res.Height = 14, 9
	res.timer = 360
	res.die = false
	return res
}
