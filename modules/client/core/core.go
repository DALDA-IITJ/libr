package core

import "client/core/config"

type Core struct {
}

type UserMessage struct {
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type MsgCert struct {
	Sender  string    `json:"sender"`
	Msg     string    `json:"msg"`
	TS      string    `json:"ts"`
	ModCert []ModSign `json:"mod_cert"`
	Sign    string    `json:"sign"`
}

// DatabaseNode represents a database node in the network.
type DatabaseNode struct {
	IP   string
	Port string
}

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type MessageResponse struct {
	Messages []Message `json:"messages"`
}

// Moderator represents a moderator node.
type Moderator struct {
	IP        string
	Port      string
	PublicKey string
}

// ModSign represents a moderator's signature response.
type ModSign struct {
	Sign      string `json:"sign"`
	PublicKey string `json:"public_key"`
}

// ModCert aggregates valid moderator signatures.
type ModCert struct {
	Msg        string
	Timestamp  string
	Signatures []ModSign // Only stores {PublicKey, Signature}
}

func InitCore() {
	config.LoadEnv() // Load .env when core initializes
}

func NewCore() *Core {
	return &Core{}
}
