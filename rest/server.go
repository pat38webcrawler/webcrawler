package rest

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

//Shutdown stop properly the webcrawler server
func (s *Server) Shutdown() {
	if s != nil {

		log.Printf("Shutting down http server")
		err := s.listener.Close()
		if err != nil {
			log.Print(err)
		}
	}
}

// RunServer starts the webcrawler server
func RunServer() error {
	httpServer, err := NewServer()
	if err != nil {
		return err
	}
	defer httpServer.Shutdown()

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-signalCh
	return nil
}
