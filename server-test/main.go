package main

import (
	"log"
	"net"

	"github.com/farinas09/go-grpc/database"
	"github.com/farinas09/go-grpc/server"
	"github.com/farinas09/go-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	server := server.NewTestServer(repo)

	s := grpc.NewServer()

	testpb.RegisterTestServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
