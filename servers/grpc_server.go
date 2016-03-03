package servers

import (
	"github.com/sayden/gubsub/grpc"
	"golang.org/x/net/context"
	"net"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/vendor/src/google.golang.org/grpc"
	"github.com/sayden/gubsub/types"
	"fmt"
"time"
	"github.com/sayden/gubsub/dispatcher"
)

// server is used to implement helloworld.GreeterServer.
type rpc_server struct{
	Port int
}

var Server rpc_server

func init(){
	Server = rpc_server{
		Port:5123,
	}
}


//NewMessage is the implementation to receive a new message across the cluster
func (s *rpc_server) NewMessage(ctx context.Context, in *grpcservice.GubsubMessage) (*grpcservice.GubsubReply, error) {
	m := types.Message{
		Data:&in.M,
		Topic:&in.T,
		Timestamp:time.Now(),
	}

	dispatcher.DispatchMessage(&m)
	return &grpcservice.GubsubReply{StatusCode:0}, nil
}

func (s *rpc_server) StartServer(port int){
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Error("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcservice.RegisterMessageServiceServer(s, &rpc_server{})
	s.Serve(lis)
}

func (s *rpc_server) SendMessage(m types.Message) (int, error){
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", s.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := grpcservice.NewMessageServiceClient(conn)

	r, err := c.NewMessage(context.Background(), grpcservice.GubsubMessage{
		M:*m.Data,
		T:*m.Topic,
	})

	if err != nil {
		return err
	}

	return r, nil
}