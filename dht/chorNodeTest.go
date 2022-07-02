package chord

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"
)

func JoinTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n := NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.getSuccessor())
	fmt.Println(n.getSuccessor())
	fmt.Println(n.Info.NodeID)
}

func SleepTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n := NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	time.Sleep(25 * time.Second)
	n.stabilize()
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	n.Stop()

}

func StabilizeTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var n *Node
	var i *NodeInfo = nil
	if ipSucc != "" {
		h.Reset()
		h.Write([]byte(ip + ":" + port))
		val = h.Sum(nil)
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n = NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	time.Sleep(5 * time.Second)
	d, err := json.Marshal("Set Pepe")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = n.Set("Pepe", "", d)
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(2 * time.Second)
	s, err := n.GetByName("Pepe")
	if err != nil {
		fmt.Println(err.Error())
	}
	var f string
	json.Unmarshal(s, &f)
	fmt.Println(f)
	n.Stop()
}

func FireTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var n *Node
	var i *NodeInfo = nil
	if ipSucc != "" {
		h.Reset()
		h.Write([]byte(ip + ":" + port))
		val = h.Sum(nil)
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n = NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	time.Sleep(60 * time.Second)
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	n.Stop()
}

func SelfTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var n *Node
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n = NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	time.Sleep(60 * time.Second)
	n.Stop()
}
