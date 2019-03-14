package emulutil

import (
	"fmt"
	"time"
)

// Storage provides ability to store and access any kind of data
type Storage map[string]interface{}

// Process simulates a separate machine with network access and multiple work functions running
type Process struct {
	Pid      int32
	network  *Network
	Contexts map[string]Storage
}

func newProcess(id int32, n *Network) *Process {
	p := Process{Pid: id, network: n}
	return &p
}

func (p *Process) start() {
	go func() {
		for {
			msg := <-p.network.in[p.Pid]
			mTarget, _ := msg.GetString()
			go p.network.workFunctions[mTarget](p, &msg)
		}
	}()
}

// Send performs a non-blocking send
func (p *Process) Send(to int32, msg Message) {
	msg.from = p.Pid
	msg.to = to
	msg.sendTime = time.Now()
	go func() { p.network.pool <- msg }()
}

// GetNeighbors provides list of the process's neighbors
func (p *Process) GetNeighbors() []int32 {
	var ns []int32
	for id := range p.network.conn[p.Pid] {
		ns = append(ns, id)
	}
	return ns
}

// Log prints string with process's id
func (p *Process) Log(s string, args ...interface{}) {
	p.network.logs <- fmt.Sprintf("[%d] ", p.Pid) + fmt.Sprintf(s, args...)
}
