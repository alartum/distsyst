package main

import (
	"time"

	"github.com/alartum/distsyst/emulutil"
)

func sayHello(p *emulutil.Process, m *emulutil.Message) {
	ns := p.GetNeighbors()
	p.Log("Hello, my neighbors: %v\n", ns)
}

func bully(p *emulutil.Process, m *emulutil.Message) {
	s, _ := m.GetString()
	msg := emulutil.Message{}
	switch s {
	case "INIT":
		msg.PutInt32(1)
	}
}

func main() {
	var conf emulutil.Config
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

	net := emulutil.NewNetwork(&conf)
	net.Launch()

	msg := emulutil.Message{}
	msg.PutString("HELLO")
	for _, id := range net.GetIds() {
		net.Initialize(id, msg)
	}

	time.Sleep(10 * time.Second)
}
