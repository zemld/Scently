package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	pb "github.com/zemld/Scently/generated/proto/perfume-hub"
	"github.com/zemld/Scently/perfume-hub/api/grpc_handlers"
	"github.com/zemld/Scently/perfume-hub/api/handlers"
	"github.com/zemld/Scently/perfume-hub/api/middleware"
	"github.com/zemld/Scently/perfume-hub/internal/db/core"
	"google.golang.org/grpc"
)

func main() {
	core.Initiate()
	defer core.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Printf("Failed to listen: %v", err)
			return
		}

		s := grpc.NewServer()
		defer s.GracefulStop()
		pb.RegisterPerfumeStorageServer(s, grpc_handlers.NewPerfumeStorageServer(core.Select, core.Update))

		if err := s.Serve(listener); err != nil {
			log.Printf("Failed to serve: %v", err)
			return
		}
	}()

	go func() {
		defer wg.Done()

		r := http.NewServeMux()

		r.Handle("/v1/perfumes/get", middleware.Auth(http.HandlerFunc(handlers.Select)))
		r.Handle("/v1/perfumes/update", middleware.Auth(http.HandlerFunc(handlers.Update)))

		http.ListenAndServe(":8000", r)
	}()

	wg.Wait()
}
