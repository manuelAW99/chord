package chord

import (
	"bytes"
)

type Address struct {
	IP   string
	Port string
}

type Agent struct {
	Name string
}

// Covert an address into a string
func getAddr(addr Address) string {
	return addr.IP + ":" + addr.Port
}

func between(a, b, key []byte) bool {
	switch bytes.Compare(a, b) {
	case -1:
		return bytes.Compare(a, key) == -1 && bytes.Compare(key, b) == -1
	case 0:
		return !bytes.Equal(a, key)
	case 1:
		return bytes.Compare(a, key) == -1 || bytes.Compare(key, b) == -1
	}
	return false
}

func betweenRightInlcude(a, b, key []byte) bool {
	return between(a, b, key) || bytes.Equal(key, b)
}
