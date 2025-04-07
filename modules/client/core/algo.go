package core

import "client/core/blockchain"

func fetchDatabaseNodes(RelevantTxs []blockchain.Transaction, timestamp string) []DatabaseNode {

	return []DatabaseNode{
		{"192.168.1.20", "8080"},
		{"192.168.1.21", "8081"},
		{"192.168.1.22", "8082"},
	}
}

// [publicKey1, publicKey2]
