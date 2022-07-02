package chord

import (
	"math/big"
)

// fingerEntry represents a single finger table entry
type fingerEntry struct {
	ID       []byte    // ID hash of (n + 2^k) mod (2^m)
	SuccNode *NodeInfo // ptr to node Succ(Node_ID)
}

// Just the FT
type fingerTable struct {
	Table []*fingerEntry
}

// newFingerEntry creates a new fingerEntry
func newFingerEntry(id []byte, nodeInfo *NodeInfo) *fingerEntry {
	return &fingerEntry{
		ID:       id,
		SuccNode: nodeInfo,
	}
}

// Calculate (n + 2^k) mod (2^m)
func calculateFingerEntryID(nodeID []byte, k, m int) []byte {
	// Clone nodeID in id
	// id = n
	id := (&big.Int{}).SetBytes(nodeID)

	// Get the offset
	// Offset = 2^k
	two := big.NewInt(2)
	offset := big.Int{}
	offset.Exp(two, big.NewInt(int64(k)), nil)

	// Sum
	// Sum = n + 2^k
	sum := big.Int{}
	sum.Add(id, &offset)

	// Get the ceiling
	// Ceil = 2^m
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	// Apply the mod
	// Id = id mod (2^m)
	id.Mod(&sum, &ceil)

	// Convert id into []byte
	return id.Bytes()
}

// Create a new FT
func newFingerTable(nodeInfo *NodeInfo, m int) fingerTable {
	// Create the table
	ft := make([]*fingerEntry, m)

	// Fill the table
	for i := range ft {
		ft[i] = newFingerEntry(calculateFingerEntryID(nodeInfo.NodeID, i, m), nodeInfo)
	}

	return fingerTable{Table: ft}
}
