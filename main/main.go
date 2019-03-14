package main

import (
	"time"

	"github.com/alartum/distsyst/emul"
)

func sayHello(p *emul.Process, m *emul.Message) {
	ns := p.GetNeighbors()
	p.Log("Hello, my neighbors: %v\n", ns)
}

func bully(p *emul.Process, m *emul.Message) {
	s, _ := m.GetString()
	msg := emul.Message{}
	switch s {
	case "INIT":
		msg.PutInt32(1)
	}
}

func main() {
	var conf emul.Config
	conf.AddEdgeUndirected(0, 4, 2)
	conf.AddEdgeUndirected(0, 5, 4)
	conf.AddEdgeUndirected(1, 5, 7)
	conf.AddEdgeUndirected(2, 1, 3)
	conf.AddEdgeUndirected(2, 3, 2)
	conf.AddEdgeUndirected(4, 1, 5)
	conf.AddEdgeUndirected(4, 5, 2)
	conf.AddEdgeUndirected(5, 2, 2)
	conf.AddEdgeUndirected(5, 3, 6)

	conf.AddWorkFunction("HELLO", sayHello)

	net := emul.NewNetwork(&conf)
	net.Launch()

	for _, id := range net.GetIds() {
		msg := new(emul.Message)
		msg.PutString("HELLO")
		net.Initialize(id, msg)
	}
	time.Sleep(1 * time.Second)
	net.Drop(0)
	time.Sleep(10 * time.Second)
}
