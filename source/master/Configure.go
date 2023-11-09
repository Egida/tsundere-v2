package master

import (
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"tsundere/source"
	"tsundere/source/database"
	"tsundere/source/master/commands/impls"
)

func Configure() error {
	server := &Server{
		Config: &ServerConfig{
			Address:    ":7700",
			PrivateKey: "resources/ssh.ppk",
		},
		server: &ssh.ServerConfig{
			ServerVersion: "SSH-2.0-Tsundere", // I should make this random or get the latest
			NoClientAuth:  true,               // We will make a login via. our terminal interface (or not)
		},
	}

	// Unmarshal our server configuration
	if err := source.Parser.Unmarshal(&server.Config, "masters"); err != nil {
		return err
	}

	content, err := os.ReadFile(server.Config.PrivateKey)
	if err != nil {
		return err
	}

	private, err := ssh.ParsePrivateKey(content)
	if err != nil {
		return err
	}

	log.Printf("Private key %s has been parsed properly", server.Config.PrivateKey)
	server.server.AddHostKey(private)

	if !server.Config.CustomLogin {
		server.server.NoClientAuth = false
		server.server.PasswordCallback = func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, database.VerifyCredentials(conn.User(), string(password))
		}
	}

	// command impls
	impls.Init()

	return server.Serve()
}
