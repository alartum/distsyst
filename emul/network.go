package emul

import (
	"fmt"
	"time"
)

// Network layer each process supposed to register in
type Network struct {
	errRate       float32
	conn          Graph
	in            map[int32]chan Message
	pool          chan Message
	workFunctions map[string]func(*Process, *Message)

	logs        chan string
	stopRequest chan bool
	IsStopped   chan bool
}

// NewNetwork constructs network by configuration given in config
func NewNetwork(config *Config) *Network {
	n := new(Network)

	n.logs = make(chan string)
	n.pool = make(chan Message)
	n.conn = make(Graph)
	n.in = make(map[int32]chan Message)
	n.workFunctions = make(map[string]func(*Process, *Message))

	for id, vs := range config.conn {
		newMap := make(map[int32]float32)
		for v, delay := range vs {
			newMap[v] = delay
		}
		n.conn[id] = newMap
		n.in[id] = make(chan Message)
		p := newProcess(id, n)
		p.start()
	}
	for name, f := range config.workFunctions {
		n.workFunctions[name] = f
	}
	return n
}

// GetIds returns ids of all nodes present in the network
func (n *Network) GetIds() []int32 {
	var ids []int32
	for id := range n.conn {
		ids = append(ids, id)
	}
	return ids
}

// Stop gives the network a signal to stop inter-process communication
func (n *Network) Stop() {
	n.stopRequest <- true
}

// Launch gives the network a signal to start inter-process communication
func (n *Network) Launch() {
	fmt.Println("[Network] Launched")
	go func() {
		for {
			select {
			case mLog := <-n.logs:
				fmt.Printf("%s", mLog)
			case mNew := <-n.pool:
				n.send(mNew)
			case mStop := <-n.stopRequest:
				if mStop == true {
					fmt.Println("[Network] Stopped")
					return
				}
			}
		}
	}()
}

func (n *Network) send(msg Message) {
	ch, ok := n.in[msg.to]
	if !ok {
		panic(ok)
	}
	var delay float32
	if msg.from != -1 {
		vs, ok := n.conn[msg.from]
		if !ok {
			panic(ok)
		}
		delay, ok = vs[msg.to]
		if !ok {
			panic(ok)
		}

	}
	go func() {
		time.Sleep(time.Duration(delay) * time.Second)
		msg.deliveryTime = time.Now()
		ch <- msg
	}()
}

// Initialize sends a start message to the given process
func (n *Network) Initialize(to int32, msg Message) {
	msg.from = -1
	msg.to = to
	msg.sendTime = time.Now()
	go func() { n.pool <- msg }()
}
