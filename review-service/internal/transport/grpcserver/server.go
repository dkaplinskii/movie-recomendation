package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"studentgit.kata.academy/movie-recommendation-platform-20/review-service/pkg/reviewpb"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewGRPCServer(port int, reviewHandler reviewpb.ReviewServiceServer) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen tcp %s: %w", address, err)
	}

	grpcServer := grpc.NewServer()

	reviewpb.RegisterReviewServiceServer(grpcServer, reviewHandler)

	return &Server{
		grpcServer: grpcServer,
		listener:   listener,
	}, nil
}

func (s *Server) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	serverErrCh := make(chan error, 1)

	go func() {
		err := s.grpcServer.Serve(s.listener)

		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serverErrCh <- fmt.Errorf("serve grpc: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := s.stop(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown grpc server: %w", err)
		}

	case err := <-serverErrCh:
		return fmt.Errorf("grpc server failed: %w", err)
	}

	return nil
}

func (s *Server) stop(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		defer close(done)
		s.grpcServer.GracefulStop()
	}()

	select {
	case <-done:
		return nil

	case <-ctx.Done():
		s.grpcServer.Stop()
		<-done
		return fmt.Errorf(
			"grpc graceful shutdown timeout: %w",
			ctx.Err(),
		)
	}
}
