package master

import "golang.org/x/crypto/ssh"

// Server holds the server configuration
type Server struct {
	Config *ServerConfig
	server *ssh.ServerConfig
}

// ServerConfig holds the ssh configuration
type ServerConfig struct {
	Address     string `json:"server.address"`
	PrivateKey  string `json:"server.private_key"`
	CustomLogin bool   `json:"server.custom_login"`
}
