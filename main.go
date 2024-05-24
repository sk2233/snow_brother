package main

import (
	"gameBase/app"
	"gameBase/config"
	"gameBase/tool"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math/rand"
	_ "snowBrother/other"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixMicro())
	ebiten.SetWindowSize(768, 672)
	ebiten.SetWindowTitle("雪人兄弟")
	config.Width, config.Height = 256, 224
	config.ShowFps = true
	//config.Debug = true
	app.RunApp(NewMainApp())
}

type MainApp struct {
	*app.MapApp
	index int
}

var (
	rooms = []string{"level1.tmx", "level2.tmx", "level3.tmx", "level4.tmx", "level5.tmx", "level6.tmx",
		"level7.tmx", "boss.tmx"}
	audios = []string{"1.mp3", "3.mp3", "2.mp3"}
)

func (s *MainApp) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		s.setIndex((s.index + len(rooms) - 1) % len(rooms))
	} else if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		s.setIndex((s.index + 1) % len(rooms))
	}
	return s.SimpleApp.Update()
}

func (s *MainApp) setIndex(index int) {
	s.index = index
	s.ReplaceMap(rooms[index])
	tool.StopAll()
	tool.Loop(audios[index%3], 60)
}

func NewMainApp() *MainApp {
	res := new(MainApp)
	res.MapApp = app.NewMapApp()
	res.index = 0
	tool.Loop(audios[0], 60)
	res.PushMap(rooms[0])
	return res
}
