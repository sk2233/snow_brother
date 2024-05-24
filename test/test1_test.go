package test

import (
	"fmt"
	"gameBase/other"
	"gameBase/sprite"
	"math/rand"
	"net"
	"testing"
)

func Test11(t *testing.T) {
	fmt.Println(rand.Intn(5))
	var num other.Sprite
	res, ok := num.(*sprite.DynamicSprite)
	fmt.Println(res, ok)
}

func Test12(t *testing.T) {
	listen, _ := net.Listen("tcp", ":80")
	accept, _ := listen.Accept()
	bs := make([]byte, 1024)
	for num, _ := accept.Read(bs); num > 0; num, _ = accept.Read(bs) {
		fmt.Println(string(bs[:num]))
	}
}
