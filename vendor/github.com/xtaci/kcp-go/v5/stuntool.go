package kcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	SERVER_ADDR = "https://www.hana-sweet.top/nyaa/"
)

type Node struct {
	Id       string           `json:"id"` // Path
	Describe string           `json:"describe"`
	Endpoint string           `json:"endpoint"`
	Peer     map[string]int64 `json:"peer"`
	LastSeen int64            `json:"lastseen"`
}

func NewNode(path, endpoint string) {
	node := Node{
		Id:       path,
		Endpoint: endpoint,
	}
	json_data, err := json.Marshal(node)
	if err != nil {
		log.Printf("error : %v", err)
	}
	resp, err := http.Post(SERVER_ADDR, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}

func JoinNode(path, endpoint string) {
	node := &Node{
		Id:       path,
		Endpoint: endpoint,
	}
	json_data, err := json.Marshal(node)
	if err != nil {
		log.Printf("error : %v", err)
	}
	resp, err := http.Post(SERVER_ADDR+path, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(node)
	// fmt.Println(node)
	resp.Body.Close()
}

func GetNode(path string) *Node {
	node := &Node{}
	resp, err := http.Get(SERVER_ADDR + path)
	if err != nil {
		log.Printf("error : %v", err)
	}
	json.NewDecoder(resp.Body).Decode(node)
	// fmt.Println(node)
	resp.Body.Close()
	return node
}

func GetNodeEndpoint(path string) string {
	n := GetNode(path)
	return n.Endpoint
}

func PingPeer(path string, conn *net.UDPConn) {
	node := GetNode(path)
	for paddr := range node.Peer {
		// Debug(paddr)
		addr, err := net.ResolveUDPAddr("udp", paddr)
		if err != nil {
			log.Printf("error : %v", err)
			continue
		}
		conn.WriteTo([]byte{}, addr)
	}

}
func PingHost(path string, conn *net.UDPConn) {
	node := GetNode(path)
	addr, err := net.ResolveUDPAddr("udp", node.Endpoint)
	if err != nil {
		log.Printf("error : %v", err)
		return
	}
	conn.WriteTo([]byte{}, addr)

}

func (node *Node) PingPeer(conn *net.UDPConn) {
	// node := GetNode(path)
	for paddr := range node.Peer {
		addr, err := net.ResolveUDPAddr("udp", paddr)
		if err != nil {
			log.Printf("error : %v", err)
			continue
		}
		conn.WriteTo([]byte{}, addr)
	}
}
func (node *Node) PingHost(conn *net.UDPConn) {
	// node := GetNode(path)
	addr, err := net.ResolveUDPAddr("udp", node.Endpoint)
	if err != nil {
		log.Printf("error : %v", err)
		return
	}
	conn.WriteTo([]byte{}, addr)
}

func sleep(s int) {
	time.Sleep(time.Millisecond * time.Duration(s))
}

func GetNodeEndpointBrute(path string) string {
	r := make(chan string)
	// quit := make(chan struct{})
	go func() {
		flag := true
		for flag {
			go func() {
				r <- GetNodeEndpoint(path)
				flag = false
			}()
			time.Sleep(time.Millisecond * 250)
		}
	}()
	return <-r
}

func JoinNodeBrute(path, endpoint string) {
	c := make(chan bool)
	go func() {
		flag := true
		for flag {
			go func() {
				JoinNode(path, endpoint)
				flag = false
				c <- true
			}()
			time.Sleep(time.Millisecond * 250)
		}
	}()
	<-c
}
func NewNodeBrute(path, endpoint string) {
	c := make(chan bool)
	// quit := make(chan struct{})
	go func() {
		flag := true
		for flag {
			go func() {
				NewNode(path, endpoint)
				flag = false
				c <- true
			}()
			time.Sleep(time.Millisecond * 250)
		}
	}()
	<-c
}
