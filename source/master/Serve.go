package master

import (
	"log"
	"tsundere/packages/utilities/sshd"
	"tsundere/source/master/views"
)

// Serve will properly start the SSH server
func (s *Server) Serve() error {
	sshListener, err := sshd.ListenSSH(s.Config.Address, s.server)
	if err != nil {
		return err
	}

	log.Printf("SSH server is now listening for connections on %s", sshListener.Addr().String())

	sshListener.HandlerFunc = views.Welcome
	sshListener.Serve()
	return nil
}
