package emulutil

// Graph describes single-directed connections
type Graph map[int32]map[int32]float32

// Config provides network topology and functions configuration
type Config struct {
	conn          Graph
	workFunctions map[string]func(*Process, *Message)
}

// AddWorkFunction adds new work function to the list
func (c *Config) AddWorkFunction(name string, f func(*Process, *Message)) {
	if c.workFunctions == nil {
		c.workFunctions = map[string]func(*Process, *Message){}
	}
	c.workFunctions[name] = f
}

// AddEdgeDirected adds directed edge from one node to another and sets its delay
func (c *Config) AddEdgeDirected(from, to int32, delay float32) {
	if c.conn == nil {
		c.conn = Graph{}
	}
	if _, ok := c.conn[from]; !ok {
		c.conn[from] = make(map[int32]float32)
	}
	c.conn[from][to] = delay
}

// AddEdgeUndirected adds edge between two nodes and sets its delay
func (c *Config) AddEdgeUndirected(from, to int32, delay float32) {
	c.AddEdgeDirected(from, to, delay)
	c.AddEdgeDirected(to, from, delay)
}
