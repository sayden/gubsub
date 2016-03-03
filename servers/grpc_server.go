package servers

import (
	"fmt"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/grpc"
	"github.com/sayden/gubsub/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type rpc_server struct {
	Port int
}

var RPC rpc_server

func init() {
	RPC = rpc_server{
		Port: 5123,
	}

	go func(){
		time.Sleep(2 * time.Second)
		RPC.StartServer()
	}()
}

//NewMessage is the implementation to receive a new message across the cluster
func (s *rpc_server) NewMessage(ctx context.Context, in *grpcservice.GubsubMessage) (*grpcservice.GubsubReply, error) {
	m := types.Message{
		Data:      &in.M,
		Topic:     &in.T,
		Timestamp: time.Now(),
	}

	dispatcher.DispatchMessageLocal(&m)
	return &grpcservice.GubsubReply{StatusCode: 0}, nil
}

func (s *rpc_server) StartServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Error("Failed starting server: failed to listen:", err)
	}

	server := grpc.NewServer()
	grpcservice.RegisterMessageServiceServer(server, &rpc_server{})
	server.Serve(lis)
}

func (s *rpc_server) SendMessage(m *types.Message) (int32, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", s.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := grpcservice.NewMessageServiceClient(conn)

	r, err := c.NewMessage(context.Background(), &grpcservice.GubsubMessage{
		M: *m.Data,
		T: *m.Topic,
	})

	if err != nil {
		return -1, err
	}

	return r.StatusCode, nil
}
