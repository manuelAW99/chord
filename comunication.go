package chord

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Request structs
type EmptyRequest struct{}
type KeyRequest struct{ Key []byte }
type NodeInfoRequest struct{ NInfo *NodeInfo }
type DataKeyRequest struct {
	Key  []byte
	Data []byte
}

// Response structs
type EmptyResponse struct{}
type NodeInfoResponse struct {
	NInfo *NodeInfo
	IsNil bool
}
type DataResponse struct {
	Data []byte
}

// Alive error
type PingResquestError struct{}

func (p PingResquestError) Error() string {
	return "server off"
}

// Interface for comunication between
// Chord nodes
// All nodes have to implemet the interface
type Comunication interface {
	GetSuccessorOf(Address) (NodeInfo, error)
	GetSuccessorOfKey(Address, []byte) (NodeInfo, error)
	GetPredecessorOf(Address) (NodeInfo, error)
	Notify(Address) error
	Ping(Address) bool
	AskForAKey(Address, []byte) (string, error)
}

// Server logic for rpc

func NewRPCServer(n *Node) *rpc.Server {
	s := rpc.NewServer()
	if err := s.Register(n); err != nil {
		panic(err.Error())
	}

	return s
}

func RunServer(s *rpc.Server, addr Address, stopC chan struct{}) {
	listener, err := net.Listen("tcp", getAddr(addr))
	if err != nil {
		panic(err.Error())
	}

	go func() {
		for {
			select {
			case <-stopC:
				fmt.Println("Deteniendo servidor")
				if err = listener.Close(); err != nil {
					panic(err)
				}
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println(err)
					return
				}
				go s.ServeConn(conn)
			}
		}
	}()

}

// Exported rpc method for GetSuccessor
func (n *Node) GetSuccessorRPC(request *EmptyRequest, response *NodeInfoResponse) error {
	succ := n.getSuccessor()
	if succ == nil {
		response.NInfo = &NodeInfo{}
		response.IsNil = true
		return nil
	}

	response.NInfo = &NodeInfo{}
	response.NInfo.NodeID = succ.NodeID
	response.NInfo.EndPoint = succ.EndPoint
	response.IsNil = false
	return nil
}

// Exported rpc method for GetSuccessorOfKey
func (n *Node) GetSuccessorOfKeyRPC(resquest *KeyRequest, response *NodeInfoResponse) error {
	info := n.findSuccessorOfKey(resquest.Key)
	if info == nil {
		response.NInfo = &NodeInfo{}
		response.IsNil = true
		return nil
	}

	response.NInfo = &NodeInfo{}
	response.NInfo.NodeID = info.NodeID
	response.NInfo.EndPoint = info.EndPoint
	response.IsNil = false
	return nil
}

// Exported rpc method for getPredecessor
func (n *Node) GetPredecessorOfRPC(request *EmptyRequest, response *NodeInfoResponse) error {
	predInfo := n.getPredecessor()
	if predInfo == nil {
		response.NInfo = &NodeInfo{}
		response.IsNil = true
		return nil
	}

	response.NInfo = &NodeInfo{}
	response.NInfo.NodeID = predInfo.NodeID
	response.NInfo.EndPoint = predInfo.EndPoint
	response.IsNil = false
	return nil
}

// Exported rpc method for notify
func (n *Node) NotifyNode(request *NodeInfoRequest, response *EmptyResponse) error {
	n.notify(request.NInfo)
	return nil
}

func (n *Node) MakePingRPC(request *EmptyRequest, response *EmptyResponse) error {
	return nil
}

func (n *Node) GetResource(request *KeyRequest, response *DataResponse) error {
	n.dbMutex.RLock()
	defer n.dbMutex.RUnlock()

	data, err := n.db.GetByName(request.Key)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	response.Data = data
	return nil
}

func (n *Node) SaveResource(request *DataKeyRequest, response *EmptyResponse) error {
	n.dbMutex.Lock()
	defer n.dbMutex.Unlock()

	return n.db.Set(request.Key, request.Data)
}

func (n *Node) DeleteResource(request *KeyRequest, response *EmptyResponse) error {
	n.dbMutex.Lock()
	defer n.dbMutex.Unlock()

	return n.db.Delete(request.Key)
}

// Client logic for rpc

// Connect the client with server at address addr
// And calls GetSuccessor of node at address addr
func getSuccessorOf(addr Address) (*NodeInfo, error) {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer client.Close()

	var response *NodeInfoResponse = &NodeInfoResponse{}
	err = client.Call("Node.GetSuccessorRPC", &EmptyRequest{}, response)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if response.IsNil {
		log.Println("Response is nil")
		return nil, nil
	}

	return response.NInfo, nil
}

func getSuccessorOfKey(addr Address, key []byte) (*NodeInfo, error) {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer client.Close()

	var response *NodeInfoResponse = &NodeInfoResponse{NInfo: nil, IsNil: false}
	err = client.Call("Node.GetSuccessorOfKeyRPC", &KeyRequest{Key: key}, response)
	if err != nil {
		client.Close()
		log.Println(err.Error())
		return nil, err
	}

	if response.IsNil {
		log.Println("Response is Nil")
		return nil, nil
	}

	if response.NInfo == nil {
		fmt.Println("response nil")
		return nil, nil
	}

	return response.NInfo, nil
}

// Ask the node at address addr for his predecessor
func getPredecessorOf(addr Address) (*NodeInfo, error) {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer client.Close()

	var response *NodeInfoResponse = &NodeInfoResponse{}
	err = client.Call("Node.GetPredecessorOfRPC", &EmptyRequest{}, response)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if response.IsNil {
		log.Println("Response is nil")
		return nil, nil
	}

	return response.NInfo, nil
}

func notifyNode(addr Address, n *NodeInfo) error {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer client.Close()

	err = client.Call("Node.NotifyNode", &NodeInfoRequest{NInfo: n}, nil)
	return err
}

func ping(addr Address) error {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		log.Println(err.Error())
		return PingResquestError{}
	}
	defer client.Close()

	return nil
}

func askForAKey(addr Address, key []byte) ([]byte, error) {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var response *DataResponse = &DataResponse{}
	err = client.Call("Node.GetResource", &KeyRequest{Key: key}, response)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func sendSet(addr Address, key []byte, data []byte) error {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer client.Close()

	return client.Call("Node.SaveResource", &DataKeyRequest{Key: key, Data: data}, &EmptyResponse{})
}

func sendDelete(addr Address, key []byte) error {
	client, err := rpc.Dial("tcp", getAddr(addr))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer client.Close()

	return client.Call("Node.DeleteResource", &KeyRequest{Key: key}, &EmptyResponse{})
}
