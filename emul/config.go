package emul

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Graph describes single-directed connections
type Graph map[int32]map[int32]float32

// Config provides network topology and functions configuration
type Config struct {
	conn          Graph
	workFunctions map[string]func(*Process, *Message)
	init          []configInit
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

type configFile struct {
	Network []configConnection
	Init    []configInit
}

type configConnection struct {
	Directed bool
	From     []int32
	To       []int32
	Delay    float32
}

type configInit struct {
	To  []int32
	Msg string
}

// LoadFromFile loads config file and constructs Config instance accordingly
func (c *Config) LoadFromFile(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[Config] Loading from \"%s\"\n", filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()
	var config configFile
	err = json.Unmarshal([]byte(byteValue), &config)
	if err != nil {
		fmt.Println(err)
	}

	for _, conn := range config.Network {
		for _, from := range conn.From {
			for _, to := range conn.To {
				if conn.Directed {
					c.AddEdgeDirected(from, to, conn.Delay)
				} else {
					c.AddEdgeUndirected(from, to, conn.Delay)
				}
			}
		}
	}

	c.init = config.Init
}
