package grpcserver

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	address    string
	log        *zap.SugaredLogger
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(port int, log *zap.SugaredLogger) (*Server, error) {
	address := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen tcp %s: %w", address, err)
	}

	grpcServer := grpc.NewServer()

	return &Server{
		address:    address,
		log:        log,
		grpcServer: grpcServer,
		listener:   listener,
	}, nil
}

func (s *Server) Start() error {
	s.log.Infof("grpc server starting: address=%s", s.address)

	if err := s.grpcServer.Serve(s.listener); err != nil {
		if err == grpc.ErrServerStopped {
			return nil
		}

		return fmt.Errorf("serve grpc: %w", err)
	}

	return nil
}

func (s *Server) Stop(timeout time.Duration) {
	done := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-done:
		s.log.Info("grpc server stopped gracefully")

	case <-timer.C:
		s.log.Warnf("grpc graceful stop timeout exceeded, forcing stop: timeout=%s", timeout)
		s.grpcServer.Stop()
		<-done
	}
}

func (s *Server) GRPCServer() *grpc.Server {
	return s.grpcServer
}
