package emul

import (
	"fmt"
	"time"
)

const fromNobody = -1

// Network layer each process supposed to register in
type Network struct {
	conn          Graph
	in            map[int32]chan *Message
	pool          chan *Message
	delivery      chan *Message
	workFunctions map[string]func(*Process, *Message)

	logs        chan string
	dropRequest chan int32
	addRequest  chan int32
	stopRequest chan bool
	IsStopped   chan bool
}

// NewNetwork constructs network by configuration given in config
func NewNetwork(config *Config) *Network {
	n := new(Network)

	n.logs = make(chan string)
	n.pool = make(chan *Message)
	n.delivery = make(chan *Message)
	n.addRequest = make(chan int32)
	n.dropRequest = make(chan int32)
	n.stopRequest = make(chan bool)
	n.conn = make(Graph)
	n.in = make(map[int32]chan *Message)
	n.workFunctions = make(map[string]func(*Process, *Message))

	for id, vs := range config.conn {
		newMap := make(map[int32]float32)
		for v, delay := range vs {
			newMap[v] = delay
		}
		n.conn[id] = newMap
		n.in[id] = make(chan *Message)
		p := newProcess(id, n)
		p.launch()
	}
	for name, f := range config.workFunctions {
		n.workFunctions[name] = f
	}

	for _, init := range config.init {
		for _, to := range init.To {
			msg := new(Message)
			msg.PutString(init.Msg)
			n.Initialize(to, msg)
		}
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

// Drop allows to drop given node from the network
func (n *Network) Drop(id int32) {
	n.dropRequest <- id
}

// Add allows to add a node with given id to the network
func (n *Network) Add(id int32) {
	n.addRequest <- id
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
			case mDeliver := <-n.delivery:
				n.deliver(mDeliver)
			case mDrop := <-n.dropRequest:
				n.dropNode(mDrop)
			case mAdd := <-n.addRequest:
				n.addNode(mAdd)
			case mStop := <-n.stopRequest:
				if mStop == true {
					fmt.Println("[Network] Stopped")
					return
				}
			}
		}
	}()
}

func (n *Network) addNode(id int32) {
	_, ok := n.in[id]
	if ok {
		fmt.Printf("[Network] Node [%d] is present, can't add\n", id)
	} else {
		n.conn[id] = make(map[int32]float32)
		n.in[id] = make(chan *Message)
		p := newProcess(id, n)
		p.launch()
	}
}

func (n *Network) dropNode(id int32) {
	ch, ok := n.in[id]
	if !ok {
		fmt.Printf("[Network] Node [%d] is not present, can't drop\n", id)
	} else {
		fmt.Printf("[Network] Dropping node [%d]\n", id)
		close(ch)
		delete(n.in, id)
		delete(n.conn, id)
		for _, ns := range n.conn {
			for u := range ns {
				delete(ns, u)
			}
		}
	}
}

func (n *Network) deliver(msg *Message) {
	msg.deliveryTime = time.Now()
	ch, ok := n.in[msg.to]
	if !ok {
		fmt.Printf("[Network] Delivery failed, host [%d] is not present\n", msg.to)
	} else {
		ch <- msg
	}
}

func (n *Network) send(msg *Message) {
	var delay float32
	if msg.from != fromNobody {
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
		n.delivery <- msg
	}()
}

// Initialize sends a start message to the given process
func (n *Network) Initialize(to int32, msg *Message) {
	msg.from = fromNobody
	msg.to = to
	msg.sendTime = time.Now()
	go func() { n.pool <- msg }()
}
