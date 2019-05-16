package main

import (
	"time"

	"github.com/alartum/distsyst/emul"
)

func sayHello(p *emul.Process, m *emul.Message) {
	ns := p.GetNeighbors()
	p.Log("Hello, my neighbors: %v\n", ns)
}

func main() {
	var conf emul.Config

	conf.LoadFromFile("config.json")
	conf.AddWorkFunction("HELLO", sayHello)

	net := emul.NewNetwork(&conf)
	net.Launch()

	time.Sleep(1 * time.Second)
	net.Drop(0)
	// time.Sleep(10 * time.Second)
}
