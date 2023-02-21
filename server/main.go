package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	pb "simpleGRPC-Quynhlx/proto/message"
	"time"
)

const socket string = "localhost:50051"
const gwSocket string = "localhost:51151"

var messageList []*pb.MessageOne

type Server struct {
	pb.MessageServiceServer
}

func main() {
	lisn, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatalln("Errored while Listen to : ", socket, err)
	}
	log.Println("Listening at ", socket)
	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &Server{})
	go s.Serve(lisn)
	time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatalln("Errored while Serving : ", socket, err)
	}

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		socket,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	// MuxServe objct
	gwmux := runtime.NewServeMux()
	// Register MessageService handler with the client connection and Mux object
	err = pb.RegisterMessageServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	// Create HTTP Server using the socket and the Mux handler
	gwServer := &http.Server{
		Addr:    gwSocket,
		Handler: gwmux,
	}
	// Make the server to Listen and Serve
	log.Println("Serving gRPC-Gateway on ", gwSocket)
	log.Fatalln(gwServer.ListenAndServe())

}

func (s *Server) GetMessage(ctx context.Context, req *pb.MessageID) (*pb.MessageOne, error) {
	log.Println("Hitted GetMessage with the message ID", req.Id)
	for _, m := range messageList {
		if m.Id == req.Id {
			return m, nil
		}
	}
	return nil, status.Errorf(
		codes.NotFound,
		"Given Message ID is not found",
	)

}

func (s *Server) CreateMessage(ctx context.Context, msg *pb.MessageOne) (*pb.MessageID, error) {

	log.Println("Hitted CreateMessage with the msg id ", msg.Id)
	log.Println("Checking whether the given msg id is already there")
	for _, m := range messageList {
		if m.Id == msg.Id {
			log.Println("The Given message ID already present. So skipping the create process")
			return nil, status.Errorf(
				codes.AlreadyExists,
				"Given message ID already exists",
			)
		}
	}
	messageList = append(messageList, msg)
	empID := pb.MessageID{Id: msg.Id}
	log.Println("Given message ID doesn't exists. Succesfully created.")
	return &empID, nil
}
