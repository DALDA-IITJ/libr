package core

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"sort"

	"github.com/DALDA-IITJ/libr/modules/core/blockchain"
)

// func fetchDatabaseNodes(RelevantTxs []blockchain.Transaction, timestamp string, k int) []DatabaseNode {

//		return []DatabaseNode{
//			{"192.168.1.20", "8080"},
//			{"192.168.1.21", "8081"},
//			{"192.168.1.22", "8082"},
//		}
//	}
//

// Step 1: Preprocessing — Compute current DB nodes
func GetActiveDBNodes(relevantTxs []blockchain.Transaction) []string {
	nodeStatus := make(map[string]bool)

	for _, tx := range relevantTxs {
		if tx.Data.Type == "DB_JOINED" {
			nodeStatus[tx.Sender] = true
		} else if tx.Data.Type == "DB_LEFT" && tx.Data.Leaver != "" {
			nodeStatus[tx.Data.Leaver] = false
		}
	}

	var activeNodes []string
	for node, isActive := range nodeStatus {
		if isActive {
			activeNodes = append(activeNodes, node)
		}
	}

	sort.Strings(activeNodes) // ensure deterministic order
	return activeNodes
}

// Step 2: Deterministically select K nodes using timestamp
func SelectKNodesFromTimestamp(timestamp string, activeNodes []string, k int) []string {
	if len(activeNodes) == 0 || k == 0 {
		return []string{}
	}

	hash := sha256.Sum256([]byte(timestamp))
	randSeed := binary.BigEndian.Uint64(hash[:8])
	randGen := rand.New(rand.NewSource(int64(randSeed)))

	selected := make(map[int]bool)
	var selectedNodes []string

	for len(selectedNodes) < k && len(selectedNodes) < len(activeNodes) {
		idx := randGen.Intn(len(activeNodes))
		if !selected[idx] {
			selected[idx] = true
			selectedNodes = append(selectedNodes, activeNodes[idx])
		}
	}

	return selectedNodes
}

// Step 3: Return list of selected DatabaseNodes
func fetchDatabaseNodes(relevantTxs []blockchain.Transaction, timestamp string, k int) []DatabaseNode {
	activeNodes := GetActiveDBNodes(relevantTxs)
	selectedPublicKeys := SelectKNodesFromTimestamp(timestamp, activeNodes, k)

	// Map from public key → (IP, Port)
	nodeMeta := make(map[string]DatabaseNode)
	for _, tx := range relevantTxs {
		if tx.Data.Type == "DB_JOINED" {
			rawJSON, ok := tx.Data.Metadata["metadata"]
			if !ok {
				continue
			}

			var parsed map[string]string
			if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
				continue
			}

			nodeMeta[tx.Sender] = DatabaseNode{
				IP:   parsed["ip"],
				Port: parsed["port"],
			}
		}
	}

	var selectedNodes []DatabaseNode
	for _, pubKey := range selectedPublicKeys {
		if node, exists := nodeMeta[pubKey]; exists {
			selectedNodes = append(selectedNodes, node)
		}
	}

	return selectedNodes
}
